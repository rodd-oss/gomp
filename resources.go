package tinyrpg

import (
	"image"
	"image/png"
	"os"

	"github.com/markbates/pkger"
)

type Frames struct {
	Frames []image.Image
	image.Config
}

func LoadResources() (map[string]Frames, error) {
	images := map[string]image.Image{}
	cfgs := map[string]image.Config{}
	sprites := map[string]Frames{}

	prefix := "/resources/sprites"
	err := pkger.Walk(prefix, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() == false {
			filename := prefix + "/" + info.Name()
			file, err := pkger.Open(filename)
			if err != nil {
				return err
			}
			defer file.Close()

			img, err := png.Decode(file)
			if err != nil {
				return err
			}

			fileCfg, err := pkger.Open(filename)
			if err != nil {
				return err
			}
			defer fileCfg.Close()

			cfg, err := png.DecodeConfig(fileCfg)
			if err != nil {
				return err
			}

			images[info.Name()] = img
			cfgs[info.Name()] = cfg
		}

		return nil
	})

	if err != nil {
		return sprites, err
	}

	sprites["big_demon_idle"] = Frames{
		Frames: []image.Image{
			images["big_demon_idle_anim_f0.png"],
			images["big_demon_idle_anim_f1.png"],
			images["big_demon_idle_anim_f2.png"],
			images["big_demon_idle_anim_f3.png"],
		},
		Config: cfgs["big_demon_idle_anim_f0.png"],
	}
	sprites["big_demon_run"] = Frames{
		Frames: []image.Image{
			images["big_demon_run_anim_f0.png"],
			images["big_demon_run_anim_f1.png"],
			images["big_demon_run_anim_f2.png"],
			images["big_demon_run_anim_f3.png"],
		},
		Config: cfgs["big_demon_run_anim_f0.png"],
	}

	sprites["big_zombie_idle"] = Frames{
		Frames: []image.Image{
			images["big_zombie_idle_anim_f0.png"],
			images["big_zombie_idle_anim_f1.png"],
			images["big_zombie_idle_anim_f2.png"],
			images["big_zombie_idle_anim_f3.png"],
		},
		Config: cfgs["big_zombie_idle_anim_f0.png"],
	}
	sprites["big_zombie_run"] = Frames{
		Frames: []image.Image{
			images["big_zombie_run_anim_f0.png"],
			images["big_zombie_run_anim_f1.png"],
			images["big_zombie_run_anim_f2.png"],
			images["big_zombie_run_anim_f3.png"],
		},
		Config: cfgs["big_zombie_run_anim_f0.png"],
	}

	sprites["elf_f_idle"] = Frames{
		Frames: []image.Image{
			images["elf_f_idle_anim_f0.png"],
			images["elf_f_idle_anim_f1.png"],
			images["elf_f_idle_anim_f2.png"],
			images["elf_f_idle_anim_f3.png"],
		},
		Config: cfgs["elf_f_idle_anim_f0.png"],
	}
	sprites["elf_f_run"] = Frames{
		Frames: []image.Image{
			images["elf_f_run_anim_f0.png"],
			images["elf_f_run_anim_f1.png"],
			images["elf_f_run_anim_f2.png"],
			images["elf_f_run_anim_f3.png"],
		},
		Config: cfgs["elf_f_run_anim_f0.png"],
	}
	sprites["floor_1"] = Frames{
		Frames: []image.Image{images["floor_1.png"]},
		Config: cfgs["floor_1.png"],
	}
	sprites["floor_2"] = Frames{
		Frames: []image.Image{images["floor_2.png"]},
		Config: cfgs["floor_2.png"],
	}
	sprites["floor_3"] = Frames{
		Frames: []image.Image{images["floor_3.png"]},
		Config: cfgs["floor_3.png"],
	}
	sprites["floor_4"] = Frames{
		Frames: []image.Image{images["floor_4.png"]},
		Config: cfgs["floor_4.png"],
	}
	sprites["floor_5"] = Frames{
		Frames: []image.Image{images["floor_5.png"]},
		Config: cfgs["floor_5.png"],
	}
	sprites["floor_6"] = Frames{
		Frames: []image.Image{images["floor_6.png"]},
		Config: cfgs["floor_6.png"],
	}
	sprites["floor_7"] = Frames{
		Frames: []image.Image{images["floor_7.png"]},
		Config: cfgs["floor_7.png"],
	}
	sprites["floor_8"] = Frames{
		Frames: []image.Image{images["floor_8.png"]},
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
