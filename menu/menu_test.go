package menu

import (
	"os"
	"reflect"
	"testing"

	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"

	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

func Test_WarpToQuickMenu(t *testing.T) {
	m := Init(&video.Video{})

	t.Run("Starts with a single scene if no game is running", func(t *testing.T) {
		got := len(menu.stack)
		want := 1
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Starts on the tabs scene if no game is running", func(t *testing.T) {
		got := menu.stack[0].Entry().label
		want := "Ludo"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	state.CoreRunning = true
	m.WarpToQuickMenu()

	t.Run("Warps at the quick menu if a game is launched", func(t *testing.T) {
		got := len(menu.stack)
		want := 3
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Warps at the quick menu if a game is launched", func(t *testing.T) {
		got := menu.stack[len(menu.stack)-1].Entry().label
		want := "Quick Menu"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("No tweens are left after warping to quick menu", func(t *testing.T) {
		got := len(menu.tweens)
		want := 0
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_fastForwardTweens(t *testing.T) {
	foo := float32(5)
	menu.tweens[&foo] = gween.New(foo, 25, 0.015, ease.OutSine)
	bar := float32(-10)
	menu.tweens[&bar] = gween.New(bar, 0, 0.9, ease.OutSine)

	menu.tweens.FastForward()

	t.Run("No tweens are left", func(t *testing.T) {
		got := len(menu.tweens)
		want := 0
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Tweened vars have the right value", func(t *testing.T) {
		got := foo
		want := float32(25)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("Tweened vars have the right value", func(t *testing.T) {
		got := bar
		want := float32(0)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func Test_buildExplorer(t *testing.T) {

	menu.stack = []Scene{}

	exec := 0
	cbMock := func(str string) {
		exec++
	}

	dirActionMock := &entry{
		label: "<Scan this directory>",
		icon:  "scan",
	}

	tmp := os.TempDir() + "/Test_buildExplorer/"

	os.RemoveAll(tmp)
	os.Mkdir(tmp, 0777)
	defer os.RemoveAll(tmp)
	os.Create(tmp + "File 1.txt")
	os.Create(tmp + "File 2.img")
	os.Create(tmp + "File 3.txt")
	os.Create(tmp + "File 4.img") // this entry will be pushed at the end by sort+prettify
	os.Create(tmp + "File 5.txt")
	os.Create(tmp + "File 6.txt")
	os.Create(tmp + "File 7.txt")
	os.Mkdir(tmp+"Folder 1", 0777)

	prettify := func(in string) string {
		if in == "File 4" {
			return "IMAGE 4"
		}
		return in
	}

	scene := buildExplorer(os.TempDir()+"/Test_buildExplorer/", []string{".img"}, cbMock, dirActionMock, prettify)
	menu.Push(scene)

	children := scene.Entry().children

	t.Run("Should display the right number of menu entries", func(t *testing.T) {
		// txt files are ignored because we filter on .img
		if len(children) != 5 {
			t.Errorf("buildExplorer = %v, want %v", len(children), 5)
		}
	})

	t.Run("Inserts the directory Action as first entry", func(t *testing.T) {
		if children[0].label != dirActionMock.label {
			t.Errorf("buildExplorer = %v, want %v", children[0].label, dirActionMock.label)
		}
	})

	t.Run("Inserts the directory .. as second entry", func(t *testing.T) {
		if children[1].label != ".." {
			t.Errorf("buildExplorer = %v, want %v", children[1].label, "..")
		}
	})

	t.Run("Files have file icon", func(t *testing.T) {
		if children[2].icon != "file" {
			t.Errorf("buildExplorer = %v, want %v", children[2].icon, "file")
		}
	})

	t.Run("Targeted files have OK callbacks", func(t *testing.T) {
		children[2].callbackOK()
		if exec != 1 {
			t.Errorf("buildExplorer = %v, want %v", exec, 1)
		}
	})

	t.Run("Folders have folder icon", func(t *testing.T) {
		if children[3].icon != "folder" {
			t.Errorf("buildExplorer = %v, want %v", children[4].icon, "folder")
		}
	})

	t.Run("Folder callback opens the folder", func(t *testing.T) {
		if len(menu.stack) != 1 {
			t.Errorf("buildExplorer = %v, want %v", len(menu.stack), 1)
		}
		children[3].callbackOK()
		if len(menu.stack) != 2 {
			t.Errorf("buildExplorer = %v, want %v", len(menu.stack), 2)
		}
	})

	t.Run("Prettifier should work", func(t *testing.T) {
		want := "IMAGE 4"
		if children[4].label != want {
			t.Errorf("buildExplorer = %v, want %v", children[3].label, want)
		}
	})
}

func TestExtractTags(t *testing.T) {
	var empty []string
	tests := []struct {
		name  string
		args  string
		want1 interface{}
		want2 interface{}
	}{
		{
			name:  "No tags",
			args:  "My Awesome Game",
			want1: "My Awesome Game",
			want2: empty,
		},
		{
			name:  "One tag",
			args:  "My Awesome Game (France)",
			want1: "My Awesome Game",
			want2: []string{"France"},
		},
		{
			name:  "Multiple tags",
			args:  "My Awesome Game (France) (v1.0)",
			want1: "My Awesome Game",
			want2: []string{"France", "v1.0"},
		},
		{
			name:  "Nested tags",
			args:  "My Awesome Game (Europe) (Fr,De,En)",
			want1: "My Awesome Game",
			want2: []string{"Europe", "Fr", "De", "En"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := extractTags(tt.args)
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
