package resources

import (
	"embed"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Frames struct {
	Frames []*ebiten.Image
	image.Config
}

//go:embed **/*.png
var fs embed.FS

func walkDir(prefix string, fn func(path string, info os.FileInfo, err error) error) error {
	dirEntries, err := fs.ReadDir(prefix)
	if err != nil {
		return err
	}

	for _, entry := range dirEntries {
		info, err := entry.Info()
		if err != nil {
			return err
		}

		err = fn(entry.Name(), info, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func Load() (frames map[string]Frames, err error) {
	images := map[string]image.Image{}
	cfgs := map[string]image.Config{}
	sprites := map[string]Frames{}

	prefix := "sprites"

	err = walkDir(prefix, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		filename := prefix + "/" + info.Name()
		file, err := fs.Open(filename)
		if err != nil {
			fmt.Println("Error opening file")
			return err
		}
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			fmt.Println("Error decoding file")
			return err
		}

		fileCfg, err := fs.Open(filename)
		if err != nil {
			fmt.Println("Error decoding file")
			return err
		}
		defer fileCfg.Close()

		cfg, err := png.DecodeConfig(fileCfg)
		if err != nil {
			fmt.Println("Error decoding cfg file")
			return err
		}

		images[info.Name()] = img
		cfgs[info.Name()] = cfg

		return nil
	})

	if err != nil {
		return nil, err
	}

	sprites["big_demon_idle"] = Frames{
		Frames: []*ebiten.Image{
			ebiten.NewImageFromImage(images["big_demon_idle_anim_f0.png"]),
			ebiten.NewImageFromImage(images["big_demon_idle_anim_f1.png"]),
			ebiten.NewImageFromImage(images["big_demon_idle_anim_f2.png"]),
			ebiten.NewImageFromImage(images["big_demon_idle_anim_f3.png"]),
		},
		Config: cfgs["big_demon_idle_anim_f0.png"],
	}
	sprites["big_demon_run"] = Frames{
		Frames: []*ebiten.Image{
			ebiten.NewImageFromImage(images["big_demon_run_anim_f0.png"]),
			ebiten.NewImageFromImage(images["big_demon_run_anim_f1.png"]),
			ebiten.NewImageFromImage(images["big_demon_run_anim_f2.png"]),
			ebiten.NewImageFromImage(images["big_demon_run_anim_f3.png"]),
		},
		Config: cfgs["big_demon_run_anim_f0.png"],
	}

	sprites["big_zombie_idle"] = Frames{
		Frames: []*ebiten.Image{
			ebiten.NewImageFromImage(images["big_zombie_idle_anim_f0.png"]),
			ebiten.NewImageFromImage(images["big_zombie_idle_anim_f1.png"]),
			ebiten.NewImageFromImage(images["big_zombie_idle_anim_f2.png"]),
			ebiten.NewImageFromImage(images["big_zombie_idle_anim_f3.png"]),
		},
		Config: cfgs["big_zombie_idle_anim_f0.png"],
	}
	sprites["big_zombie_run"] = Frames{
		Frames: []*ebiten.Image{
			ebiten.NewImageFromImage(images["big_zombie_run_anim_f0.png"]),
			ebiten.NewImageFromImage(images["big_zombie_run_anim_f1.png"]),
			ebiten.NewImageFromImage(images["big_zombie_run_anim_f2.png"]),
			ebiten.NewImageFromImage(images["big_zombie_run_anim_f3.png"]),
		},
		Config: cfgs["big_zombie_run_anim_f0.png"],
	}

	sprites["elf_f_idle"] = Frames{
		Frames: []*ebiten.Image{
			ebiten.NewImageFromImage(images["elf_f_idle_anim_f0.png"]),
			ebiten.NewImageFromImage(images["elf_f_idle_anim_f1.png"]),
			ebiten.NewImageFromImage(images["elf_f_idle_anim_f2.png"]),
			ebiten.NewImageFromImage(images["elf_f_idle_anim_f3.png"]),
		},
		Config: cfgs["elf_f_idle_anim_f0.png"],
	}
	sprites["elf_f_run"] = Frames{
		Frames: []*ebiten.Image{
			ebiten.NewImageFromImage(images["elf_f_run_anim_f0.png"]),
			ebiten.NewImageFromImage(images["elf_f_run_anim_f1.png"]),
			ebiten.NewImageFromImage(images["elf_f_run_anim_f2.png"]),
			ebiten.NewImageFromImage(images["elf_f_run_anim_f3.png"]),
		},
		Config: cfgs["elf_f_run_anim_f0.png"],
	}
	sprites["floor_1"] = Frames{
		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_1.png"])},
		Config: cfgs["floor_1.png"],
	}
	sprites["floor_2"] = Frames{
		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_2.png"])},
		Config: cfgs["floor_2.png"],
	}
	sprites["floor_3"] = Frames{
		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_3.png"])},
		Config: cfgs["floor_3.png"],
	}
	sprites["floor_4"] = Frames{
		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_4.png"])},
		Config: cfgs["floor_4.png"],
	}
	sprites["floor_5"] = Frames{
		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_5.png"])},
		Config: cfgs["floor_5.png"],
	}
	sprites["floor_6"] = Frames{
		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_6.png"])},
		Config: cfgs["floor_6.png"],
	}
	sprites["floor_7"] = Frames{
		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_7.png"])},
		Config: cfgs["floor_7.png"],
	}
	sprites["floor_8"] = Frames{
		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_8.png"])},
		Config: cfgs["floor_8.png"],
	}

	return sprites, nil
}

func LoadLevel() [][]string {
	a := "floor_1"
	b := "floor_2"
	c := "floor_3"
	d := "floor_4"

	level := [][]string{
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, b, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, c, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, c, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, d, a, a, a, a, a, a, a},
		{a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a, a},
	}

	return level
}
