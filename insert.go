package sqlutils

import "fmt"
import "reflect"
import "strings"
import "strconv"
import "time"

func BuildInsertColumnsFromMap(cols map[string]interface{}) (string, []interface{}) {
	var columnNames []string
	var valueFields []string
	var values []interface{}
	var i int = 1

	for col, arg := range cols {
		columnNames = append(columnNames, col)
		valueFields = append(valueFields, fmt.Sprintf("$%d", i))
		values = append(values, arg)
		i++
	}
	sql := "( " + strings.Join(columnNames, ", ") + " )" +
		" VALUES ( " + strings.Join(valueFields, ", ") + " )"
	return sql, values
}

// Build Insert Column Clause from a struct type.
func BuildInsertColumns(val interface{}) (string, []interface{}) {
	t := reflect.ValueOf(val).Elem()
	typeOfT := t.Type()

	var columnNames []string
	var valueFields []string
	var values []interface{}
	var fieldId int = 1

	for i := 0; i < t.NumField(); i++ {
		var fieldType = typeOfT.Field(i)
		var tag reflect.StructTag = fieldType.Tag
		var field reflect.Value = t.Field(i)

		var columnName *string = GetColumnNameFromTag(&tag)
		if columnName == nil {
			continue
		}

		var attributes = GetColumnAttributesFromTag(&tag)

		// if it's a serial column (with auto-increment, we can simply skip)
		if _, ok := attributes["serial"]; ok {
			continue
		}

		var val interface{} = field.Interface()

		// if time is null or with zero value, just skip it.
		if fieldType.Type.String() == "*time.Time" {
			if timeVal, ok := val.(*time.Time); ok {
				if timeVal == nil || timeVal.Unix() == -62135596800 {
					continue
				}
			}
		}

		if attributes["date"] {
			switch val.(type) {
			case string:
				if val == "" {
					continue
				}
			}
		}

		columnNames = append(columnNames, *columnName)
		valueFields = append(valueFields, "$"+strconv.Itoa(fieldId))
		values = append(values, val)
		fieldId++
	}
	var values2 []string
	for _, e := range values {
		if _, ok := e.(time.Time); ok {
			//			cc := e.(time.Time)
			//			timeNow := cc.Format("2006-01-02 15:04:05")
			values2 = append(values2, e.(time.Time).String())
		} else if _, ok := e.(int64); ok {
			sst := strconv.FormatInt(e.(int64), 10)
			values2 = append(values2, sst)
		} else if _, ok := e.(int); ok {
			sst2 := strconv.Itoa(e.(int))
			values2 = append(values2, sst2)
		} else {
			values2 = append(values2, e.(string))
		}
	}
	//	return "( " + strings.Join(columnNames, ", ") + " ) " +
	//		"VALUES ( " + strings.Join(valueFields, ", ") + " )", values
	return "( " + strings.Join(columnNames, ", ") + " ) " +
		"VALUES ( '" + strings.Join(values2, "','") + "' )", values
}

func BuildInsertColumnsMulti(val interface{}) (string, []interface{}) {
	//	t := reflect.ValueOf(val).Elem()
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	var retsql string
	kai := 0
	kai2 := 0
	for i := 0; i < l; i++ {
		t := reflect.ValueOf(ret[i])
		typeOfT := t.Type()

		var columnNames []string
		var valueFields []string
		var values []interface{}
		var fieldId int = 1

		for i := 0; i < t.NumField(); i++ {
			var fieldType = typeOfT.Field(i)
			var tag reflect.StructTag = fieldType.Tag
			var field reflect.Value = t.Field(i)

			var columnName *string = GetColumnNameFromTag(&tag)
			if columnName == nil {
				continue
			}

			var attributes = GetColumnAttributesFromTag(&tag)

			// if it's a serial column (with auto-increment, we can simply skip)
			if _, ok := attributes["serial"]; ok {
				continue
			}

			var val interface{} = field.Interface()

			// if time is null or with zero value, just skip it.
			if fieldType.Type.String() == "*time.Time" {
				if timeVal, ok := val.(*time.Time); ok {
					if timeVal == nil || timeVal.Unix() == -62135596800 {
						continue
					}
				}
			}

			if attributes["date"] {
				switch val.(type) {
				case string:
					if val == "" {
						continue
					}
				}
			}
			if kai == 0 {

				columnNames = append(columnNames, *columnName)
			}
			valueFields = append(valueFields, "$"+strconv.Itoa(fieldId))
			values = append(values, val)
			fieldId++
		}
		if kai == 0 {
			retsql = "( " + strings.Join(columnNames, ", ") + " ) " + "VALUES ( '"
			kai++
		}
		//		fmt.Println(columnNames, valueFields, ret, "-----------------")
		var values2 []string
		for _, e := range values {
			if _, ok := e.(time.Time); ok {
				//                      cc := e.(time.Time)
				//                      timeNow := cc.Format("2006-01-02 15:04:05")
				values2 = append(values2, e.(time.Time).String())
			} else if _, ok := e.(int64); ok {
				sst := strconv.FormatInt(e.(int64), 10)
				values2 = append(values2, sst)
			} else if _, ok := e.(int); ok {
				sst2 := strconv.Itoa(e.(int))
				values2 = append(values2, sst2)
			} else {
				values2 = append(values2, e.(string))
			}
		}
		if kai2 == 0 {
			retsql += strings.Join(values2, "','") + "' )"
			kai2++
			continue
		} else {
			retsql += ",('" + strings.Join(values2, "','") + "' )"

		}
	}
	//      return "( " + strings.Join(columnNames, ", ") + " ) " +
	//              "VALUES ( " + strings.Join(valueFields, ", ") + " )", values
	//	return "( " + strings.Join(columnNames, ", ") + " ) " +
	//		"VALUES ( '" + strings.Join(values2, "','") + "' )", values
	//retsql = "( " + strings.Join(columnNames, ", ") + " ) " + "VALUES ( '"
	return retsql, ret
}

func BuildInsertClause(val interface{}) (string, []interface{}) {
	tableName := GetTableName(val)
	fmt.Println(tableName)
	tableName = tableName[0 : len(tableName)-1]
	sql, values := BuildInsertColumns(val)
	return "INSERT INTO " + tableName + sql, values
}
func BuildInsertClause2(val interface{}) (string, []interface{}) {
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	tableName := GetTableName(v.Index(0).Interface())
	fmt.Println(tableName, "xxxxxx")
	tableName = tableName[0 : len(tableName)-1]
	sql, values := BuildInsertColumnsMulti(val)
	return "INSERT INTO " + tableName + sql, values
}
