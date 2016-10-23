package model

import (
	"net/http"
	"strconv"

	"github.com/antonholmquist/jason"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	sq "github.com/lann/squirrel"
  
	"../util"

)

var (
	brandSelect = `
		SELECT json_build_object(
  	  	'id', CAST(brand.id AS text), 
 		'name', brand.name,
  		'brandSynonym', (
			SELECT COALESCE(array_to_json(array_agg(json_build_object(
			    	'id', CAST(brand_synonym.id AS text),
		      		'name', brand_synonym.name
       			 ))), '{}'::json)
        		FROM brand_synonym
			WHERE brand.id = brand_synonym.brand_id  
      			)
		) 
		AS brand
		FROM brand `
)

type BrandResult struct {
	Brand types.JSONText `db:"brand"`
}

type Brand struct{}

//Handle query request for brand api
func (r *Brand) HandleQuery(w http.ResponseWriter, req *http.Request) {
	conn := GetConnection()
  	writer := util.NewJsonWriter(w)
	results, err := conn.Queryx(brandSelect)

	if err != nil {
		panic(err)
	}

	result := BrandResult{}
	writer.RootObject(func() {
		writer.KeyValue("status", "success")
		writer.Array("content", func() {
			for results.Next() {
				err := results.StructScan(&result)
				if err != nil {
					panic(err)
				}
				writer.RawValue([]byte(result.Brand))
			}
		})
	})
}

//Handle read request for brand api
func (r *Brand) HandleRead(w http.ResponseWriter, req *http.Request) {
	conn := GetConnection()
 	writer := util.NewJsonWriter(w)
	result := BrandResult{}
	id := req.URL.Query().Get(":id")
	sql := brandSelect + "WHERE id = $1;"
	err := conn.Get(&result, sql, id)
	if err != nil {
		panic(err)
	}
	writer.RootObject(func() {
		writer.KeyValue("status", "success")
		writer.Array("content", func() {
			writer.RawValue([]byte(result.Brand))
		})
	})
}

//Handle create request for brand api
func (r *Brand) HandleCreate(w http.ResponseWriter, req *http.Request) {
 	conn := GetConnection()
  	body := util.ParseBody(req.Body)
	data := util.ParseJson(body)
	content := util.ParseJsonContent(data)

	brandStatement := sq.Insert("brand").PlaceholderFormat(sq.Dollar)
	brandStatement = brandStatement.Columns("id", "name")
	synonymCount := 0
	synonymStatement := sq.Insert("brand_synonym").PlaceholderFormat(sq.Dollar)
	synonymStatement = synonymStatement.Columns("id", "brand_id", "name")
	for _, brand := range content {
		brandValues := parseBrandValues(brand)
		synonyms := parseSynonymArray(brand)
		for _, synonym := range synonyms {
			synonymValues := parseSynonymValues(synonym, brand)
			synonymStatement = synonymStatement.Values(synonymValues...)
			synonymCount++
		}
		brandStatement = brandStatement.Values(brandValues...)
	}

	tx := conn.MustBegin()
	brandSql, brandArgs, brandErr := brandStatement.ToSql()

	if brandErr != nil {
		panic(brandErr)
	}
	synonymSql, synonymArgs, synonymErr := synonymStatement.ToSql()
	tx.MustExec(brandSql, brandArgs...)
	if synonymCount > 0 {
		if synonymErr != nil {
			panic(synonymErr)
		}
		tx.MustExec(synonymSql, synonymArgs...)
	}
	tx.Commit()

	writer := util.NewJsonWriter(w)
	writer.RootObject(func() {
		writer.KeyValue("status", "success")
	})
}

//Handle update request for brand api
func (r *Brand) HandleUpdate(w http.ResponseWriter, req *http.Request) {
	conn := GetConnection()
  	body := util.ParseBody(req.Body)
	data := util.ParseJson(body)
	content := util.ParseJsonContent(data)

	for _, brand := range content {
		brandStatement := sq.Update("brand").PlaceholderFormat(sq.Dollar)
		brandValues := parseBrandValues(brand)
		brandStatement = brandStatement.Where(sq.Eq{"id": brandValues[0]})
		brandStatement = brandStatement.Set("name", brandValues[1])
		brandSql, brandArgs, brandErr := brandStatement.ToSql()
		if brandErr != nil {
			panic(brandErr)
		}

		synonyms := parseSynonymArray(brand)
		tx := conn.MustBegin()
		for _, synonym := range synonyms {
			synonymStatement := sq.Update("brand_synonym").PlaceholderFormat(sq.Dollar)
			synonymValues := parseSynonymValues(synonym, brand)
			synonymStatement = synonymStatement.Where(sq.Eq{"id": synonymValues[0]})
			synonymStatement = synonymStatement.Set("brand_id", synonymValues[1])
			synonymStatement = synonymStatement.Set("name", synonymValues[2])
			synonymSql, synonymArgs, synonymErr := synonymStatement.ToSql()
			if synonymErr != nil {
				panic(synonymErr)
			}
			tx.MustExec(synonymSql, synonymArgs...)
		}
		tx.MustExec(brandSql, brandArgs...)
		tx.Commit()
	}

	writer := util.NewJsonWriter(w)
	writer.RootObject(func() {
		writer.KeyValue("status", "success")
	})
}

//Parse brand json into string array
func parseBrandValues(brand *jason.Object) []interface{} {
	var values []interface{}
	id, idErr := brand.GetString("id")
	if idErr != nil {
		panic(idErr)
	}
	parsedID, parseIDErr := strconv.ParseInt(id, 10, 64)
	if parseIDErr != nil {
		panic(parseIDErr)
	}
	name, nameErr := brand.GetString("name")
	if nameErr != nil {
		panic(nameErr)
	}
	values = append(values, parsedID)
	values = append(values, name)
	return values
}

//Parse brand synonym json into string array
func parseSynonymValues(synonym *jason.Object, brand *jason.Object) []interface{} {
	var values []interface{}
	id, idErr := synonym.GetString("id")
	if idErr != nil {
		panic(idErr)
	}
	parsedID, parseIDErr := strconv.ParseInt(id, 10, 64)
	if parseIDErr != nil {
		panic(parseIDErr)
	}
	brandID, brandIDErr := brand.GetString("brandID")
	if brandIDErr != nil {
		panic(brandIDErr)
	}
	parsedBrandID, parseBrandIDErr := strconv.ParseInt(brandID, 10, 64)
	if parseBrandIDErr != nil {
		panic(parseBrandIDErr)
	}
	name, nameErr := synonym.GetString("name")
	if nameErr != nil {
		panic(nameErr)
	}
	values = append(values, parsedID)
	values = append(values, parsedBrandID)
	values = append(values, name)
	return values
}

//Parse brand synonym json collection into json array
func parseSynonymArray(brand *jason.Object) []*jason.Object {
	synonyms, err := brand.GetObjectArray("brandSynonym")
	if err != nil {
		return []*jason.Object{}
	}
	return synonyms
}
