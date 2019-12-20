package data

// UpdateOptions
type UpdateOptions struct {
	ArrayFilters             *ArrayFilters // A set of filters specifying to which array elements an update should apply
	BypassDocumentValidation *bool         // If true, allows the write to opt-out of document level validation
	Upsert                   *bool         // When true, creates a new document if no document matches the query
}

// Update returns a pointer to a new UpdateOptions
func Update() *UpdateOptions {
	return &UpdateOptions{}
}

// SetArrayFilters specifies a set of filters specifying to which array elements an update should apply
func (uo *UpdateOptions) SetArrayFilters(af ArrayFilters) *UpdateOptions {
	uo.ArrayFilters = &af
	return uo
}

// SetBypassDocumentValidation allows the write to opt-out of document level validation.
func (uo *UpdateOptions) SetBypassDocumentValidation(b bool) *UpdateOptions {
	uo.BypassDocumentValidation = &b
	return uo
}

// SetUpsert allows the creation of a new document if not document matches the query
func (uo *UpdateOptions) SetUpsert(b bool) *UpdateOptions {
	uo.Upsert = &b
	return uo
}

// MergeUpdateOptions combines the argued UpdateOptions into a single UpdateOptions in a last-one-wins fashion
func MergeUpdateOptions(opts ...*UpdateOptions) *UpdateOptions {
	uOpts := Update()
	for _, uo := range opts {
		if uo == nil {
			continue
		}
		if uo.ArrayFilters != nil {
			uOpts.ArrayFilters = uo.ArrayFilters
		}
		if uo.BypassDocumentValidation != nil {
			uOpts.BypassDocumentValidation = uo.BypassDocumentValidation
		}
		if uo.Upsert != nil {
			uOpts.Upsert = uo.Upsert
		}
	}

	return uOpts
}
