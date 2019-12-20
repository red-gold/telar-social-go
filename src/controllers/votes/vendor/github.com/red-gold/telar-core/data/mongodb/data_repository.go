package mongodb

import (
	"context"
	"fmt"

	d "github.com/red-gold/telar-core/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type DataRepositoryMongo struct {
	Client MongoDatabase
}

type DataSingleResult struct {
	result   *mongo.SingleResult
	err      error
	noResult bool
}

type DataResult struct {
	result *mongo.Cursor
	ctx    context.Context
	err    error
}

// NewDataRepositoryMongo create new data repository for mongodb.
func NewDataRepositoryMongo(client MongoDatabase) d.Repository {
	return &DataRepositoryMongo{Client: client}
}

// CreateIndex creates multiple indexes in the collection specified by the indexes.
func (m *DataRepositoryMongo) CreateIndex(collectionName string, indexes map[string]interface{}) <-chan error {
	result := make(chan error)
	go func() {

		var (
			err        error
			collection *mongo.Collection
			ctx        context.Context
		)

		collection, err = m.Client.GetCollection(collectionName)
		if err != nil {
			fmt.Printf("Get collection %s err (%s)! \n", collectionName, err.Error())
			result <- err
		}

		ctx, err = m.Client.GetContext()
		if err != nil {
			fmt.Printf("Get context err (%s)! \n", err.Error())
			result <- err
		}

		var indexList []mongo.IndexModel

		for key, value := range indexes {
			indexOption := &options.IndexOptions{}
			indexOption = indexOption.SetBackground(true)
			index := mongo.IndexModel{Keys: bson.M{key: value}, Options: indexOption}
			indexList = append(indexList, index)
		}

		_, err = collection.Indexes().CreateMany(ctx, indexList)
		result <- err
		close(result)
	}()

	return result
}

// Save storing the data object.
func (m *DataRepositoryMongo) Save(collectionName string, data interface{}) <-chan d.RepositoryResult {
	result := make(chan d.RepositoryResult)
	go func() {

		var (
			err        error
			collection *mongo.Collection
			ctx        context.Context
		)

		collection, err = m.Client.GetCollection(collectionName)
		if err != nil {
			fmt.Printf("Get collection %s err (%s)! \n", collectionName, err.Error())
			result <- d.RepositoryResult{Error: err}
		}

		ctx, err = m.Client.GetContext()
		if err != nil {
			fmt.Printf("Get context err (%s)! \n", err.Error())
			result <- d.RepositoryResult{Error: err}
		}

		_, err = collection.InsertOne(ctx, &data)
		result <- d.RepositoryResult{Error: err}
		close(result)
	}()

	return result
}

// Save storing a list of objects.
func (m *DataRepositoryMongo) SaveMany(collectionName string, data []interface{}) <-chan d.RepositoryResult {
	result := make(chan d.RepositoryResult)
	go func() {

		var (
			err        error
			collection *mongo.Collection
			ctx        context.Context
		)

		collection, err = m.Client.GetCollection(collectionName)
		if err != nil {
			fmt.Printf("Get collection %s err (%s)! \n", collectionName, err.Error())
			result <- d.RepositoryResult{Error: err}
		}

		ctx, err = m.Client.GetContext()
		if err != nil {
			fmt.Printf("Get context err (%s)! \n", err.Error())
			result <- d.RepositoryResult{Error: err}
		}
		insertOption := &options.InsertManyOptions{}
		insertOption = insertOption.SetOrdered(false)
		_, err = collection.InsertMany(ctx, data, insertOption)
		result <- d.RepositoryResult{Error: err}
		close(result)
	}()

	return result
}

// Find get list of object.
func (m *DataRepositoryMongo) Find(collectionName string, filter interface{}, limit int64, skip int64, sort map[string]int) <-chan d.QueryResult {
	result := make(chan d.QueryResult)
	go func() {

		var (
			err        error
			collection *mongo.Collection
			ctx        context.Context
			cur        *mongo.Cursor
		)

		collection, err = m.Client.GetCollection(collectionName)
		if err != nil {
			fmt.Printf("Get collection %s err (%s)! \n", collectionName, err.Error())
			result <- &DataResult{err: err}
		}

		ctx, err = m.Client.GetContext()
		if err != nil {
			fmt.Printf("Get context err (%s)! \n", err.Error())
			result <- &DataResult{err: err}
		}
		findOptions := options.Find()
		if len(sort) > 0 {
			findOptions.SetSort(sort)
		}
		if skip > 0 {
			findOptions.SetSkip(skip)
		}
		if limit > 0 {
			findOptions.SetLimit(limit)
		}

		// Execute query
		cur, err = collection.Find(ctx, filter, findOptions)
		if err != nil {
			fmt.Printf("Find cursor err (%s)! \n", err.Error())
			result <- &DataResult{err: err}
		}

		result <- &DataResult{result: cur}
		close(result)
	}()

	return result
}

// FindOne get object list
func (m *DataRepositoryMongo) FindOne(collectionName string, filter interface{}) <-chan d.QuerySingleResult {
	result := make(chan d.QuerySingleResult)
	go func() {

		var (
			err        error
			collection *mongo.Collection
			ctx        context.Context
			findResult *mongo.SingleResult
		)

		collection, err = m.Client.GetCollection(collectionName)
		if err != nil {
			fmt.Printf("Get collection %s err (%s)! \n", collectionName, err.Error())
			result <- &DataSingleResult{err: err}
		}

		ctx, err = m.Client.GetContext()
		if err != nil {
			fmt.Printf("Get context err (%s)! \n", err.Error())
			result <- &DataSingleResult{err: err}
		}

		// Execute query
		findResult = collection.FindOne(ctx, filter)
		err = findResult.Err()
		if err != nil {

			fmt.Printf("Find result error (%s)! \n", err.Error())
			result <- &DataSingleResult{err: err}
		}
		result <- &DataSingleResult{result: findResult}
		close(result)
	}()

	return result
}

// Update update object.
func (m *DataRepositoryMongo) Update(collectionName string, filter interface{}, data interface{}, opts ...*d.UpdateOptions) <-chan d.RepositoryResult {
	result := make(chan d.RepositoryResult)
	go func() {

		var (
			err          error
			collection   *mongo.Collection
			ctx          context.Context
			updateResult *mongo.UpdateResult
		)
		// Get Collection
		collection, err = m.Client.GetCollection(collectionName)
		if err != nil {
			fmt.Printf("Get collection %s err (%s)! \n", collectionName, err.Error())
			result <- d.RepositoryResult{Error: err}
		}

		// Get Context
		ctx, err = m.Client.GetContext()
		if err != nil {
			fmt.Printf("Get context err (%s)! \n", err.Error())
			result <- d.RepositoryResult{Error: err}
		}

		// Map options
		updateOptions := &options.UpdateOptions{}
		if opts != nil {
			mergedOpts := d.MergeUpdateOptions(opts...)
			if mergedOpts.Upsert != nil {
				updateOptions.SetUpsert(*mergedOpts.Upsert)
			}
			if mergedOpts.BypassDocumentValidation != nil {
				updateOptions.SetBypassDocumentValidation(*mergedOpts.BypassDocumentValidation)
			}
			arrayFilter := options.ArrayFilters{}
			if mergedOpts.ArrayFilters != nil {
				arrayFilter.Filters = mergedOpts.ArrayFilters.Filters
			}
			updateOptions.SetArrayFilters(arrayFilter)
		}

		// Execute update
		updateResult, err = collection.UpdateOne(ctx, filter, data, updateOptions)
		if err != nil {
			fmt.Printf("Update error (%s)! \n", err.Error())
			result <- d.RepositoryResult{Error: err}
		}

		// Returnt Result
		result <- d.RepositoryResult{Result: updateResult.ModifiedCount}
		close(result)
	}()

	return result
}

// UpdateMany bulk update one object.
func (m *DataRepositoryMongo) BulkUpdateOne(collectionName string, bulkData []d.BulkUpdateOne) <-chan d.RepositoryResult {
	result := make(chan d.RepositoryResult)
	go func() {

		var (
			err          error
			collection   *mongo.Collection
			ctx          context.Context
			updateResult *mongo.BulkWriteResult
		)
		collection, err = m.Client.GetCollection(collectionName)
		if err != nil {
			fmt.Printf("Get collection %s err (%s)! \n", collectionName, err.Error())
			result <- d.RepositoryResult{Error: err}
		}
		ctx, err = m.Client.GetContext()
		if err != nil {
			fmt.Printf("Get context err (%s)! \n", err.Error())
			result <- d.RepositoryResult{Error: err}
		}
		bulkOptions := &options.BulkWriteOptions{}
		bulkOptions.SetOrdered(false)
		var mongoModels []mongo.WriteModel
		for _, bulkItem := range bulkData {
			model := mongo.NewUpdateOneModel()
			model = model.SetFilter(bulkItem.Filter)
			model = model.SetUpdate(bulkItem.Data)
			mongoModels = append(mongoModels, model)
		}

		// Execute bulk update model
		updateResult, err = collection.BulkWrite(ctx, mongoModels, bulkOptions)
		if err != nil {
			fmt.Printf("Bulk Update error (%s)! \n", err.Error())
			result <- d.RepositoryResult{Error: err}
		}

		result <- d.RepositoryResult{Result: updateResult.ModifiedCount}
		close(result)
	}()

	return result
}

// Delete get object list
func (m *DataRepositoryMongo) Delete(collectionName string, filter interface{}, justOne bool) <-chan d.RepositoryResult {
	result := make(chan d.RepositoryResult)
	go func() {

		var (
			err          error
			collection   *mongo.Collection
			deleteResult *mongo.DeleteResult
			ctx          context.Context
		)
		collection, err = m.Client.GetCollection(collectionName)
		if err != nil {
			fmt.Printf("Get collection %s err (%s)! \n", collectionName, err.Error())
			result <- d.RepositoryResult{Error: err}
		}

		ctx, err = m.Client.GetContext()
		if err != nil {
			fmt.Printf("Get context err (%s)! \n", err.Error())
			result <- d.RepositoryResult{Error: err}
		}

		// Execute delete
		if justOne {
			deleteResult, err = collection.DeleteOne(ctx, filter)
		} else {
			deleteResult, err = collection.DeleteMany(ctx, filter)
		}
		if err != nil {
			fmt.Printf("Delete error (%s)! \n", err.Error())
			result <- d.RepositoryResult{Error: err}
		}

		result <- d.RepositoryResult{Result: deleteResult.DeletedCount}
		close(result)
	}()

	return result
}

// Decode single data result decoding
func (sr *DataSingleResult) Decode(v interface{}) error {
	if sr.result == nil {
		return nil
	}
	return sr.result.Decode(v)
}

// Decode single data result decoding
func (sr *DataSingleResult) NoResult() bool {
	return sr.noResult
}

// Error single data result error
func (sr *DataSingleResult) Error() error {
	sr.noResult = false
	if sr.err != nil {
		if sr.err.Error() == "mongo: no documents in result" {
			sr.noResult = true
			return nil
		}
		return sr.err
	}
	return nil

}

// Next data result iterator
func (sr *DataResult) Next() bool {
	status := sr.result.Next(sr.ctx)
	if !status {
		// defer sr.result.Close(sr.ctx)
	}
	return status
}

// Close close cursor
func (sr *DataResult) Close() {
	sr.result.Close(sr.ctx)
}

// Decode multi data result decoding
func (sr *DataResult) Decode(v interface{}) error {
	err := sr.result.Decode(v)
	if err != nil {
		fmt.Printf("Decode err (%s)! \n", err.Error())
	}
	return err
}

// Error data result error
func (sr *DataResult) Error() error {
	return sr.err
}
