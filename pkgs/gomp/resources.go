/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

// import (
// 	"embed"
// 	"image"
// 	"image/png"
// 	"log"
// 	"math/rand"
// 	"os"
// 	"strconv"

// 	"github.com/hajimehoshi/ebiten/v2"
// )

// type Sprite struct {
// 	Frames []*ebiten.Image
// 	image.Config
// }

// func walkDir(fs embed.Fs, prefix string, fn func(path string, info os.FileInfo, err error) error) error {
// 	dirEntries, err := fs.ReadDir(prefix)
// 	if err != nil {
// 		return err
// 	}

// 	for _, entry := range dirEntries {
// 		info, err := entry.Info()
// 		if err != nil {
// 			return err
// 		}

// 		err = fn(entry.Name(), info, nil)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func CreateResources() {

// }

// func Load() (sprites map[string]Sprite, err error) {
// 	images := map[string]image.Image{}
// 	cfgs := map[string]image.Config{}
// 	sprites = make(map[string]Sprite)

// 	prefix := "sprites"

// 	err = walkDir(prefix, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if info.IsDir() {
// 			return nil
// 		}

// 		filename := prefix + "/" + info.Name()
// 		file, err := fs.Open(filename)
// 		if err != nil {
// 			log.Println("Error opening file")
// 			return err
// 		}
// 		defer file.Close()

// 		img, err := png.Decode(file)
// 		if err != nil {
// 			log.Println("Error decoding file")
// 			return err
// 		}

// 		fileCfg, err := fs.Open(filename)
// 		if err != nil {
// 			log.Println("Error decoding file")
// 			return err
// 		}
// 		defer fileCfg.Close()

// 		cfg, err := png.DecodeConfig(fileCfg)
// 		if err != nil {
// 			log.Println("Error decoding cfg file")
// 			return err
// 		}

// 		images[info.Name()] = img
// 		cfgs[info.Name()] = cfg

// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	sprites["big_demon_idle"] = Sprite{
// 		Frames: []*ebiten.Image{
// 			ebiten.NewImageFromImage(images["big_demon_idle_anim_f0.png"]),
// 			ebiten.NewImageFromImage(images["big_demon_idle_anim_f1.png"]),
// 			ebiten.NewImageFromImage(images["big_demon_idle_anim_f2.png"]),
// 			ebiten.NewImageFromImage(images["big_demon_idle_anim_f3.png"]),
// 		},
// 		Config: cfgs["big_demon_idle_anim_f0.png"],
// 	}
// 	sprites["big_demon_run"] = Sprite{
// 		Frames: []*ebiten.Image{
// 			ebiten.NewImageFromImage(images["big_demon_run_anim_f0.png"]),
// 			ebiten.NewImageFromImage(images["big_demon_run_anim_f1.png"]),
// 			ebiten.NewImageFromImage(images["big_demon_run_anim_f2.png"]),
// 			ebiten.NewImageFromImage(images["big_demon_run_anim_f3.png"]),
// 		},
// 		Config: cfgs["big_demon_run_anim_f0.png"],
// 	}

// 	sprites["big_zombie_idle"] = Sprite{
// 		Frames: []*ebiten.Image{
// 			ebiten.NewImageFromImage(images["big_zombie_idle_anim_f0.png"]),
// 			ebiten.NewImageFromImage(images["big_zombie_idle_anim_f1.png"]),
// 			ebiten.NewImageFromImage(images["big_zombie_idle_anim_f2.png"]),
// 			ebiten.NewImageFromImage(images["big_zombie_idle_anim_f3.png"]),
// 		},
// 		Config: cfgs["big_zombie_idle_anim_f0.png"],
// 	}
// 	sprites["big_zombie_run"] = Sprite{
// 		Frames: []*ebiten.Image{
// 			ebiten.NewImageFromImage(images["big_zombie_run_anim_f0.png"]),
// 			ebiten.NewImageFromImage(images["big_zombie_run_anim_f1.png"]),
// 			ebiten.NewImageFromImage(images["big_zombie_run_anim_f2.png"]),
// 			ebiten.NewImageFromImage(images["big_zombie_run_anim_f3.png"]),
// 		},
// 		Config: cfgs["big_zombie_run_anim_f0.png"],
// 	}

// 	sprites["elf_f_idle"] = Sprite{
// 		Frames: []*ebiten.Image{
// 			ebiten.NewImageFromImage(images["elf_f_idle_anim_f0.png"]),
// 			ebiten.NewImageFromImage(images["elf_f_idle_anim_f1.png"]),
// 			ebiten.NewImageFromImage(images["elf_f_idle_anim_f2.png"]),
// 			ebiten.NewImageFromImage(images["elf_f_idle_anim_f3.png"]),
// 		},
// 		Config: cfgs["elf_f_idle_anim_f0.png"],
// 	}
// 	sprites["elf_f_run"] = Sprite{
// 		Frames: []*ebiten.Image{
// 			ebiten.NewImageFromImage(images["elf_f_run_anim_f0.png"]),
// 			ebiten.NewImageFromImage(images["elf_f_run_anim_f1.png"]),
// 			ebiten.NewImageFromImage(images["elf_f_run_anim_f2.png"]),
// 			ebiten.NewImageFromImage(images["elf_f_run_anim_f3.png"]),
// 		},
// 		Config: cfgs["elf_f_run_anim_f0.png"],
// 	}
// 	sprites["floor_1"] = Sprite{
// 		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_1.png"])},
// 		Config: cfgs["floor_1.png"],
// 	}
// 	sprites["floor_2"] = Sprite{
// 		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_2.png"])},
// 		Config: cfgs["floor_2.png"],
// 	}
// 	sprites["floor_3"] = Sprite{
// 		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_3.png"])},
// 		Config: cfgs["floor_3.png"],
// 	}
// 	sprites["floor_4"] = Sprite{
// 		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_4.png"])},
// 		Config: cfgs["floor_4.png"],
// 	}
// 	sprites["floor_5"] = Sprite{
// 		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_5.png"])},
// 		Config: cfgs["floor_5.png"],
// 	}
// 	sprites["floor_6"] = Sprite{
// 		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_6.png"])},
// 		Config: cfgs["floor_6.png"],
// 	}
// 	sprites["floor_7"] = Sprite{
// 		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_7.png"])},
// 		Config: cfgs["floor_7.png"],
// 	}
// 	sprites["floor_8"] = Sprite{
// 		Frames: []*ebiten.Image{ebiten.NewImageFromImage(images["floor_8.png"])},
// 		Config: cfgs["floor_8.png"],
// 	}
// 	sprites["hp"] = Sprite{
// 		Frames: []*ebiten.Image{
// 			ebiten.NewImageFromImage(images["hp1.png"]),
// 			ebiten.NewImageFromImage(images["hp2.png"]),
// 			ebiten.NewImageFromImage(images["hp3.png"]),
// 			ebiten.NewImageFromImage(images["hp4.png"]),
// 		},
// 		Config: cfgs["hp1.png"],
// 	}

// 	sprites["ability_firepizza"] = Sprite{
// 		Frames: []*ebiten.Image{
// 			ebiten.NewImageFromImage(images["fire-pizza.png"]),
// 		},
// 		Config: cfgs["fire-pizza.png"],
// 	}

// 	sprites["area"] = Sprite{
// 		Frames: []*ebiten.Image{
// 			ebiten.NewImageFromImage(images["area.png"]),
// 		},
// 		Config: cfgs["area.png"],
// 	}

// 	return sprites, nil
// }

// func LoadLevel(width int, height int) [][]string {
// 	level := make([][]string, height)

// 	// generate random level
// 	for i := 0; i < height; i++ {
// 		row := make([]string, width)

// 		for j := 0; j < width; j++ {
// 			row[j] = "floor_" + strconv.Itoa(rand.Intn(8)+1)
// 		}
// 		level[i] = row
// 	}

// 	return level
// }
