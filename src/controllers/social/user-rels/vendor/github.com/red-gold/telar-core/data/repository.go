package data

// RepositoryResult is a struct to wrap repository result
// so its easy to use it in channel
type RepositoryResult struct {
	Error  error
	Result interface{}
}

type BulkUpdateOne struct {
	Filter interface{}
	Data   interface{}
}

type ArrayFilters struct {
	Filters []interface{} // The filters to apply
}

type Repository interface {
	Save(collectionName string, data interface{}) <-chan RepositoryResult
	SaveMany(collectionName string, data []interface{}) <-chan RepositoryResult
	Find(collectionName string, filter interface{}, limit int64, skip int64, sort map[string]int) <-chan QueryResult
	FindOne(collectionName string, filter interface{}) <-chan QuerySingleResult
	Update(collectionName string, filter interface{}, data interface{}, opts ...*UpdateOptions) <-chan RepositoryResult
	BulkUpdateOne(collectionName string, bulkData []BulkUpdateOne) <-chan RepositoryResult
	Delete(collectionName string, filter interface{}, justOne bool) <-chan RepositoryResult
	CreateIndex(collectionName string, indexes map[string]interface{}) <-chan error
}

type QuerySingleResult interface {
	Decode(v interface{}) error
	NoResult() bool
	Error() error
}

type QueryResult interface {
	Next() bool
	Decode(v interface{}) error
	Error() error
	Close()
}
