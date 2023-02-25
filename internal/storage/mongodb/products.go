package storage

import (
	"context"
	"github.com/ilyadubrovsky/product-rest-api/internal/product"
	"github.com/ilyadubrovsky/product-rest-api/pkg/logging"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var productsCollection = "products"

type productsMongo struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewProductsMongo(database *mongo.Database, logger *logging.Logger) product.Repository {
	return &productsMongo{collection: database.Collection(productsCollection), logger: logger}
}

func (s *productsMongo) Create(ctx context.Context, dto product.CreateProductDTO) (string, error) {
	res, err := s.collection.InsertOne(ctx, dto)
	if err != nil {
		s.logger.Tracef("Name: %s, Type: %s, Description: %s, InStock: %d, Material: %s, Color: %s",
			dto.Name, dto.Type, dto.Description, dto.InStock,
			dto.Characteristics.Material, dto.Characteristics.Color)
		s.logger.Error("failed to InsertOne due to error: %v", err)
		return "", product.ErrInternalServer
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	s.logger.Trace(res.InsertedID)
	s.logger.Error("failed to get a hex of objectid")
	return "", product.ErrBadRequest
}

func (s *productsMongo) FindAll(ctx context.Context) ([]product.Product, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		s.logger.Errorf("failed to Find due to error: %v", err)
		return nil, product.ErrInternalServer
	}

	var prdcts []product.Product

	if err = cur.All(ctx, &prdcts); err != nil {
		s.logger.Errorf("failed to Decode elements due to error: %v", err)
		return nil, product.ErrInternalServer
	}

	if len(prdcts) == 0 {
		return nil, product.ErrNotFound
	}

	return prdcts, nil
}

func (s *productsMongo) FindOne(ctx context.Context, id string) (product.Product, error) {
	prdct := product.Product{}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Tracef("id: %s", id)
		s.logger.Errorf("failed to get a object id from id due to error: %v", err)
		return product.Product{}, product.ErrBadRequest
	}

	filter := bson.M{"_id": oid}

	res := s.collection.FindOne(ctx, filter)

	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return product.Product{}, product.ErrNotFound
	}

	if err = res.Decode(&prdct); err != nil {
		s.logger.Tracef("id: %s", id)
		s.logger.Errorf("failed to Decode a result due to error: %v", err)
		return product.Product{}, product.ErrInternalServer
	}

	return prdct, nil
}

func (s *productsMongo) FullyUpdate(ctx context.Context, id string, dto product.FullyUpdateProductDTO) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Tracef("id: %s", id)
		s.logger.Errorf("failed to get object id from id due to error: %v", err)
		return product.ErrBadRequest
	}

	productBytes, err := bson.Marshal(dto)
	if err != nil {
		s.logger.Tracef("oid: %s", oid)
		s.logger.Debugf("Name: %s, Type: %s, InStock: %d, Description: %s, Characteristics: %s",
			dto.Name, dto.Type, dto.InStock, dto.Description, dto.Characteristics)
		s.logger.Errorf("failed to bson marshal a object due to error: %v", err)
		return product.ErrInternalServer
	}

	replacement := bson.M{}
	if err = bson.Unmarshal(productBytes, &replacement); err != nil {
		s.logger.Tracef("oid: %s", oid)
		s.logger.Errorf("failed to bson unmarshal a bytes of object to bson structure due to error: %v", err)
		return product.ErrInternalServer
	}

	filter := bson.M{"_id": oid}

	res := s.collection.FindOneAndReplace(ctx, filter, replacement)

	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return product.ErrNotFound
	}

	return nil
}

func (s *productsMongo) PartiallyUpdate(ctx context.Context, id string, dto product.PartiallyUpdateProductDTO) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Tracef("id: %s", id)
		s.logger.Errorf("failed to get a object id from id due to error: %v", err)
		return product.ErrBadRequest
	}

	prdctBytes, err := bson.Marshal(dto)
	if err != nil {
		s.logger.Tracef("oid: %s", oid)
		s.logger.Errorf("failed to bson marshal a objet due to error: %v", err)
		return product.ErrInternalServer
	}

	prdctBSON := bson.M{}
	if err = bson.Unmarshal(prdctBytes, &prdctBSON); err != nil {
		s.logger.Tracef("oid: %s", oid)
		s.logger.Errorf("failed to bson unmarshal a bytes of objects due to error: %v", err)
		return product.ErrInternalServer
	}

	filter := bson.M{"_id": oid}

	update := bson.M{"$set": prdctBSON}

	res := s.collection.FindOneAndUpdate(ctx, filter, update)
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return product.ErrNotFound
	}

	return nil
}

func (s *productsMongo) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Tracef("id: %s", id)
		s.logger.Errorf("failed to get a object id from id due to error: %v", err)
		return product.ErrBadRequest
	}

	filter := bson.M{"_id": oid}

	res := s.collection.FindOneAndDelete(ctx, filter)

	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return product.ErrNotFound
	}

	if err != nil {
		s.logger.Tracef("oid: %s", oid)
		s.logger.Errorf("failed to DeleteOne due to error: %v", err)
		return product.ErrInternalServer
	}

	return nil
}
