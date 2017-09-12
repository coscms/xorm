// 兼容golang模板引擎中的逻辑符号，与其保持一致

package builder

import "fmt"

// Le defines <= condition
type Le Lte

var _ Cond = Le{}

// WriteTo write SQL to Writer
func (le Le) WriteTo(w Writer) error {
	return WriteMap(w, le, "<=")
}

// And implements And with other conditions
func (le Le) And(conds ...Cond) Cond {
	return And(le, And(conds...))
}

// Or implements Or with other conditions
func (le Le) Or(conds ...Cond) Cond {
	return Or(le, Or(conds...))
}

// IsValid tests if this Eq is valid
func (le Le) IsValid() bool {
	return len(le) > 0
}

// Ge defines >= condition
type Ge map[string]interface{}

var _ Cond = Ge{}

// WriteTo write SQL to Writer
func (ge Ge) WriteTo(w Writer) error {
	return WriteMap(w, ge, ">=")
}

// And implements And with other conditions
func (ge Ge) And(conds ...Cond) Cond {
	return And(ge, And(conds...))
}

// Or implements Or with other conditions
func (ge Ge) Or(conds ...Cond) Cond {
	return Or(ge, Or(conds...))
}

// IsValid tests if this Eq is valid
func (ge Ge) IsValid() bool {
	return len(ge) > 0
}

// Ne defines not equal conditions
type Ne map[string]interface{}

var _ Cond = Ne{}

// WriteTo writes SQL to Writer
func (ne Ne) WriteTo(w Writer) error {
	var args = make([]interface{}, 0, len(ne))
	var i = 0
	for k, v := range ne {
		switch v.(type) {
		case []int, []int64, []string, []int32, []int16, []int8:
			if err := NotIn(k, v).WriteTo(w); err != nil {
				return err
			}
		case expr:
			if _, err := fmt.Fprintf(w, "%s<>(", w.Key(k)); err != nil {
				return err
			}

			if err := v.(expr).WriteTo(w); err != nil {
				return err
			}

			if _, err := fmt.Fprintf(w, ")"); err != nil {
				return err
			}
		case *Builder:
			if _, err := fmt.Fprintf(w, "%s<>(", w.Key(k)); err != nil {
				return err
			}

			if err := v.(*Builder).WriteTo(w); err != nil {
				return err
			}

			if _, err := fmt.Fprintf(w, ")"); err != nil {
				return err
			}
		default:
			if _, err := fmt.Fprintf(w, "%s<>?", w.Key(k)); err != nil {
				return err
			}
			args = append(args, v)
		}
		if i != len(ne)-1 {
			if _, err := fmt.Fprint(w, " AND "); err != nil {
				return err
			}
		}
		i = i + 1
	}
	w.Append(args...)
	return nil
}

// And implements And with other conditions
func (ne Ne) And(conds ...Cond) Cond {
	return And(ne, And(conds...))
}

// Or implements Or with other conditions
func (ne Ne) Or(conds ...Cond) Cond {
	return Or(ne, Or(conds...))
}

// IsValid tests if this condition is valid
func (ne Ne) IsValid() bool {
	return len(ne) > 0
}
