package util

import "testing"

func TestHashedpassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hashedpassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hashedpassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Hashedpassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
