package orm

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

// GatewayPool represents a gateway pool
type GatewayPool struct {
	// URL is the connection string
	URL string
	// Isolation for each gateway instance creates a new schema and set the
	// search path to this schema
	Isolation bool

	m  map[string]*Gateway
	mu sync.RWMutex
}

// ReadDir loads all script commands from a given directory. Note that all
// scripts should have .sql extension and support the database driver.
func (p *GatewayPool) ReadDir(fileSystem FileSystem, schema ...string) error {
	var errs ErrorSlice

	for _, name := range schema {
		gateway, err := p.Get(name)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		if err = gateway.ReadDir(fileSystem); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// Migrate migrates all gateway
func (p *GatewayPool) Migrate(fileSystem FileSystem, schema ...string) error {
	var errs ErrorSlice

	for _, name := range schema {
		gateway, err := p.Get(name)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		if err = p.schema(gateway, name); err != nil {
			errs = append(errs, p.error(name, err))
		}

		if err = gateway.Migrate(fileSystem); err != nil {
			errs = append(errs, p.error(name, err))
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// Get returns a gateway for given key
func (p *GatewayPool) Get(name string) (*Gateway, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.m == nil {
		p.m = make(map[string]*Gateway)
	}

	gateway, ok := p.m[name]
	if ok {
		return gateway, nil
	}

	addr, err := p.url(name)
	if err != nil {
		return nil, p.error(name, err)
	}

	if gateway, err = Connect(addr); err != nil {
		return nil, p.error(name, err)
	}

	p.m[name] = gateway

	return gateway, nil
}

// Close closes all gateways
func (p *GatewayPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var errs ErrorSlice

	for key, gateway := range p.m {
		if err := gateway.Close(); err != nil {
			errs = append(errs, p.error(key, err))
		}

		delete(p.m, key)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (p *GatewayPool) schema(gateway *Gateway, name string) error {
	if !p.Isolation {
		return nil
	}

	param := Map{
		"schema": name,
	}

	query := SQL("CREATE SCHEMA IF NOT EXISTS {{schema}};", param)

	_, err := gateway.Exec(query)
	return err
}

func (p *GatewayPool) url(name string) (string, error) {
	if !p.Isolation {
		return p.URL, nil
	}

	addr, err := url.Parse(p.URL)
	if err != nil {
		return "", err
	}

	if addr.Scheme == "postgres" {
		schema := []string{name, "public"}

		values := addr.Query()
		values.Set("application_name", name)
		values.Set("search_path", strings.Join(schema, ","))

		addr.RawQuery = values.Encode()

		return addr.String(), nil
	}

	return "", fmt.Errorf("not supported driver %q", addr.Scheme)
}

func (p *GatewayPool) error(name string, err error) error {
	return fmt.Errorf("name: %v error: %v", name, err)
}
