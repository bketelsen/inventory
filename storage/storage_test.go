// Package storage provides an in-memory implementation of the inventory.Storage interface
// It is the only storage implementation currently available
package storage

import (
	"reflect"
	"testing"

	"github.com/bketelsen/inventory"
)

func TestMemoryStorage_StoreReport(t *testing.T) {
	type fields struct {
		reports map[string]inventory.Report
	}
	type args struct {
		report inventory.Report
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Has Hostname",
			fields: fields{
				reports: make(map[string]inventory.Report),
			},
			args: args{
				report: inventory.Report{
					Host: inventory.Host{
						HostName: "test-host",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "No Hostname",
			fields: fields{
				reports: make(map[string]inventory.Report),
			},
			args: args{
				report: inventory.Report{
					Host: inventory.Host{
						HostName: "",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemoryStorage{
				reports: tt.fields.reports,
			}
			if err := ms.StoreReport(tt.args.report); (err != nil) != tt.wantErr {
				t.Errorf("MemoryStorage.StoreReport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryStorage_GetReport(t *testing.T) {
	type fields struct {
		reports map[string]inventory.Report
	}
	type args struct {
		hostname string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   inventory.Report
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemoryStorage{
				reports: tt.fields.reports,
			}
			got, got1 := ms.GetReport(tt.args.hostname)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryStorage.GetReport() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MemoryStorage.GetReport() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemoryStorage_GetAllReports(t *testing.T) {
	type fields struct {
		reports map[string]inventory.Report
	}
	tests := []struct {
		name   string
		fields fields
		want   []inventory.Report
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemoryStorage{
				reports: tt.fields.reports,
			}
			if got := ms.GetAllReports(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryStorage.GetAllReports() = %v, want %v", got, tt.want)
			}
		})
	}
}
