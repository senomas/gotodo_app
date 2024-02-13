package sqlite

import service "github.com/senomas/gotodo_service"

type TodoFilter struct {
	categoryID  FilterInt
	done        FilterBool
	title       FilterString
	description FilterString
	category    FilterString
}

// Generate implements service.TodoFilter.
func (f *TodoFilter) Generate(qryWhere service.QueryBuilder) {
	f.title.Generate(qryWhere)
	f.description.Generate(qryWhere)
	f.category.Generate(qryWhere)
	f.categoryID.Generate(qryWhere)
	f.done.Generate(qryWhere)
}

// Category implements service.TodoFilter.
func (f *TodoFilter) Category() service.FilterString {
	f.category.field = "category.name"
	return &f.category
}

// CategoryID implements service.TodoFilter.
func (f *TodoFilter) CategoryID() service.FilterInt {
	f.categoryID.field = "category.id"
	return &f.categoryID
}

// Description implements service.TodoFilter.
func (f *TodoFilter) Description() service.FilterString {
	f.description.field = "description"
	return &f.description
}

// Done implements service.TodoFilter.
func (f *TodoFilter) Done() service.FilterBool {
	f.done.field = "done"
	return &f.done
}

// Title implements service.TodoFilter.
func (f *TodoFilter) Title() service.FilterString {
	f.title.field = "title"
	return &f.title
}

// Filter implements service.TodoService.
func (TodoService) Filter() service.TodoFilter {
	return &TodoFilter{}
}
