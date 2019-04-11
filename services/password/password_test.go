package password

import "testing"

// Tests both HashAndSalt and ComparePlaintextWithEncrypted
func TestHashAndSaltAndComparePlaintextWithEncypted(t *testing.T) {
	validpass, _ := HashAndSalt("testpassword")
	badpass, _ := HashAndSalt("badpass")

	type args struct {
		plaintextPassword string
		encryptedPassword string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid Password",
			args: args{
				plaintextPassword: "testpassword",
				encryptedPassword: validpass,
			},
			want: true,
		},
		{
			name: "Invalid Password",
			args: args{
				plaintextPassword: "testpassword",
				encryptedPassword: badpass,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComparePlaintextWithEncypted(tt.args.plaintextPassword, tt.args.encryptedPassword); got != tt.want {
				t.Errorf("ComparePlaintextWithEncypted() = %v, want %v", got, tt.want)
			}
		})
	}
}
