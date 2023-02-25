# product-rest-api
REST API for a clothing store product service using [Gin](https://gin-gonic.com/) and [MongoDB](https://github.com/mongodb/mongo-go-driver)

# endpoints
GET /products -- get all products -- can be 200/404/500 401/403
GET /products/:id -- get product by id -- can be 200/404/500 401/403
POST /products -- create product -- can be 201/400/500 401/403
PUT /products/:id -- update produt by id -- can be 204/400/404/500 401/403
PATCH /products/:id -- partially update producy by id -- can be 204/400/404/500 401/403
DELETE /products/:id -- delete product by id -- can be 204/404/500 401/403