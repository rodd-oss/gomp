/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

<- Монтажер сука Donated 50 RUB
*/

package legacy

//import (
//	"embed"
//	"fmt"
//	"image/png"
//	"io/fs"
//	"log"
//	"os"
//)
//
//func walkDir(fs fs.ReadDirFS, prefix string, fn func(path string, info os.FileInfo, err error) error) error {
//	dirEntries, err := fs.ReadDir(prefix)
//	if err != nil {
//		return err
//	}
//
//	for _, entry := range dirEntries {
//		info, err := entry.Info()
//		if err != nil {
//			return err
//		}
//
//		err = fn(entry.Name(), info, nil)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func CreateSpriteResource(fs embed.FS, prefix string) func(filename string) SpriteData {
//	sprites := make(map[string]SpriteData)
//
//	err := walkDir(fs, prefix, func(path string, info os.FileInfo, err error) error {
//		if err != nil {
//			return err
//		}
//
//		if info.IsDir() {
//			return nil
//		}
//
//		filename := prefix + "/" + info.Name()
//		file, err := fs.Open(filename)
//		if err != nil {
//			log.Println("Error opening file")
//			return err
//		}
//		defer file.Close()
//
//		img, err := png.Decode(file)
//		if err != nil {
//			log.Println("Error decoding file")
//			return err
//		}
//
//		fileCfg, err := fs.Open(filename)
//		if err != nil {
//			log.Println("Error decoding file")
//			return err
//		}
//		defer fileCfg.Close()
//
//		cfg, err := png.DecodeConfig(fileCfg)
//		if err != nil {
//			log.Println("Error decoding cfg file")
//			return err
//		}
//
//		sprites[info.Name()] = SpriteData{
//			Image:  img,
//			Config: cfg,
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		panic(err)
//	}
//
//	return func(filename string) SpriteData {
//		if sprite, ok := sprites[filename]; ok {
//			return sprite
//		}
//
//		panic(fmt.Sprint("File <", filename, "> does not exist in <", prefix, "> resource."))
//	}
//}
