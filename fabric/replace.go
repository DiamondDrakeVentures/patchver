package fabric

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/stretchr/objx"
)

type Manifest struct {
	raw objx.Map

	depends    map[string]string
	recommends map[string]string
	suggests   map[string]string
	breaks     map[string]string
	conflicts  map[string]string
}

// Depends replaces value for key in depends field.
func (m *Manifest) Depends(key, value string) bool {
	return m.replace(m.depends, key, value)
}

// Recommends replaces value for key in recommends field.
func (m *Manifest) Recommends(key, value string) bool {
	return m.replace(m.recommends, key, value)
}

// Suggests replaces value for key in suggests field.
func (m *Manifest) Suggests(key, value string) bool {
	return m.replace(m.suggests, key, value)
}

// Breaks replaces value for key in breaks field.
func (m *Manifest) Breaks(key, value string) bool {
	return m.replace(m.breaks, key, value)
}

// Conflicts replaces value for key in conflicts field.
func (m *Manifest) Conflicts(key, value string) bool {
	return m.replace(m.conflicts, key, value)
}

// JSON encodes m back to JSON.
func (m *Manifest) JSON() (string, error) {
	return m.json()
}

// MustJSON encodes m back to JSON and panics on error.
func (m *Manifest) MustJSON() string {
	j, err := m.json()
	if err != nil {
		panic(err)
	}

	return j
}

func (m *Manifest) replace(target map[string]string, key, value string) bool {
	if _, has := target[key]; has {
		target[key] = value

		return true
	}

	return false
}

func (m *Manifest) json() (string, error) {
	if len(m.depends) > 0 {
		m.raw.Set("depends", m.depends)
	}

	if len(m.recommends) > 0 {
		m.raw.Set("recommends", m.recommends)
	}

	if len(m.suggests) > 0 {
		m.raw.Set("suggests", m.suggests)
	}

	if len(m.breaks) > 0 {
		m.raw.Set("breaks", m.breaks)
	}

	if len(m.conflicts) > 0 {
		m.raw.Set("conflicts", m.conflicts)
	}

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(m.raw)
	return buffer.String(), err
}

// FromJSON decodes Manifest from json.
func FromJSON(json string) (*Manifest, error) {
	return fromJSON(json)
}

// MustFromJSON decodes Manifest from json and panics on error.
func MustFromJSON(json string) *Manifest {
	j, err := fromJSON(json)
	if err != nil {
		panic(err)
	}

	return j
}

func fromJSON(json string) (*Manifest, error) {
	m, err := objx.FromJSON(json)
	if err != nil {
		return nil, err
	}

	man := new(Manifest)
	man.raw = m

	man.depends = make(map[string]string)
	man.recommends = make(map[string]string)
	man.suggests = make(map[string]string)
	man.breaks = make(map[string]string)
	man.conflicts = make(map[string]string)

	if m.Has("depends") {
		man.populateField(man.depends, "depends")
	}

	if m.Has("recommends") {
		man.populateField(man.recommends, "recommends")
	}

	if m.Has("suggests") {
		man.populateField(man.suggests, "suggests")
	}

	if m.Has("breaks") {
		man.populateField(man.breaks, "breaks")
	}

	if m.Has("conflicts") {
		man.populateField(man.conflicts, "conflicts")
	}

	return man, nil
}

func (m *Manifest) populateField(target map[string]string, key string) error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("cannot parse field %s: %v", key, r)
		}
	}()

	if !m.raw.Has(key) {
		return fmt.Errorf("invalid key %s", key)
	}

	for k, v := range m.raw.Get(key).MustObjxMap() {
		if s, ok := v.(string); ok {
			target[k] = s
		} else {
			err = fmt.Errorf("invalid value for key %s: %v", k, v)
			return err
		}
	}

	return nil
}
