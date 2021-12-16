package wallpapers

import "testing"

func TestCalculatePath(t *testing.T) {
	type args struct {
		Url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "a url",
			args: args{
				Url: "https://five.sh/files/wallpapers/Rain/14105517898360.jpg",
			},
			want: "d:\\wallpapers/five.sh/files/wallpapers/Rain/14105517898360.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculatePath(tt.args.Url)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculatePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculatePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
