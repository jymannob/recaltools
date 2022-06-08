package recaltools

import (
	"sync"
	"testing"
)

// Funtional testing
func TestFavBackup_Backup(t *testing.T) {

	type fields struct {
		RomsDir    []string
		Gamelists  []string
		FormatJson bool
		Verbose    bool
		RestoreBkp bool
		wg         sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Backup",
			fields{
				RomsDir:    []string{"./testdata/roms"},
				FormatJson: true,
				Verbose:    true,
				RestoreBkp: false,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := &FavBackup{
				RomsDir:    tt.fields.RomsDir,
				Gamelists:  tt.fields.Gamelists,
				FormatJson: tt.fields.FormatJson,
				Verbose:    tt.fields.Verbose,
				RestoreBkp: tt.fields.RestoreBkp,
				wg:         tt.fields.wg,
			}
			if err := fb.Backup(); (err != nil) != tt.wantErr {
				t.Errorf("FavBackup.Backup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Funtional testing
func TestFavBackup_Restore(t *testing.T) {
	type fields struct {
		RomsDir    []string
		Gamelists  []string
		FormatJson bool
		Verbose    bool
		RestoreBkp bool
		wg         sync.WaitGroup
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Restore",
			fields{
				RomsDir:    []string{"./testdata/roms"},
				FormatJson: false,
				Verbose:    true,
				RestoreBkp: true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := &FavBackup{
				RomsDir:    tt.fields.RomsDir,
				Gamelists:  tt.fields.Gamelists,
				FormatJson: tt.fields.FormatJson,
				Verbose:    tt.fields.Verbose,
				RestoreBkp: tt.fields.RestoreBkp,
				wg:         tt.fields.wg,
			}
			if err := fb.Restore(); (err != nil) != tt.wantErr {
				t.Errorf("FavBackup.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFavBackup_restoreSystem(t *testing.T) {
	type fields struct {
		RomsDir    []string
		Gamelists  []string
		FormatJson bool
		Verbose    bool
		RestoreBkp bool
		wg         sync.WaitGroup
	}
	type args struct {
		gamelist string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"Restore nes system",
			fields{
				RomsDir:    []string{"./testdata/roms"},
				FormatJson: false,
				Verbose:    true,
				RestoreBkp: true,
			},
			args{"./testdata/roms/nes/gamelist.xml"},
		},
		{
			"Restore bad xml",
			fields{
				RomsDir:    []string{"./testdata/roms"},
				FormatJson: false,
				Verbose:    true,
				RestoreBkp: true,
			},
			args{"./testdata/roms/testSystem/badFile.xml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := &FavBackup{
				RomsDir:    tt.fields.RomsDir,
				Gamelists:  tt.fields.Gamelists,
				FormatJson: tt.fields.FormatJson,
				Verbose:    tt.fields.Verbose,
				RestoreBkp: tt.fields.RestoreBkp,
				wg:         tt.fields.wg,
			}
			fb.wg.Add(1)
			fb.restoreSystem(tt.args.gamelist)
		})
	}
}

func TestFavBackup_backupSystem(t *testing.T) {
	type fields struct {
		RomsDir    []string
		Gamelists  []string
		FormatJson bool
		Verbose    bool
		RestoreBkp bool
		wg         sync.WaitGroup
	}
	type args struct {
		gamelist string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"Backup megadrive system",
			fields{
				RomsDir:    []string{"./testdata/roms"},
				FormatJson: false,
				Verbose:    true,
				RestoreBkp: false,
			},
			args{"./testdata/roms/megadrive/gamelist.xml"},
		},
		{
			"Backup bad xml",
			fields{
				RomsDir:    []string{"./testdata/roms"},
				FormatJson: false,
				Verbose:    true,
				RestoreBkp: false,
			},
			args{"./testdata/roms/testSystem/badFile.xml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := &FavBackup{
				RomsDir:    tt.fields.RomsDir,
				Gamelists:  tt.fields.Gamelists,
				FormatJson: tt.fields.FormatJson,
				Verbose:    tt.fields.Verbose,
				RestoreBkp: tt.fields.RestoreBkp,
				wg:         tt.fields.wg,
			}
			fb.wg.Add(1)
			fb.backupSystem(tt.args.gamelist)
		})
	}
}
