/*
* Copyright (C) 2020 The poly network Authors
* This file is part of The poly network library.
*
* The poly network is free software: you can redistribute it and/or modify
* it under the terms of the GNU Lesser General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* The poly network is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
* GNU Lesser General Public License for more details.
* You should have received a copy of the GNU Lesser General Public License
* along with The poly network . If not, see <http://www.gnu.org/licenses/>.
 */
package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type Config struct {
	RpcEndPoints struct {
		PolyNetwork string `yaml:"polynetwork"`
		Carbon      string `yaml:"carbon"`
	} `yaml:"rpcEndPoints"`
	BroadCaster struct {
		PolyNetwork struct {
			FileName string `yaml:"fileName"`
			Password string `yaml:"password"`
		} `yaml:"polynetwork"`
	} `yaml:"broadcaster"`
}

func ReadFile(fileName string) ([]byte, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: open file %s error %s", fileName, err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: ioutil.ReadAll %s error %s", fileName, err)
	}
	return data, nil
}

func NewConfig(configFilePath string) *Config {
	fileContent, err := ReadFile(configFilePath)
	if err != nil {
		return nil
	}
	conf := &Config{}
	err = yaml.Unmarshal(fileContent, conf)
	if err != nil {
		return nil
	}

	return conf
}
