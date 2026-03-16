/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package main

import (
	"fmt"
)

// Point представляет точку с координатами X и Y.
type Point struct {
	X, Y float64
}

// Polygon представляет полигон как набор точек.
type Polygon []Point

// Projection представляет проекцию полигона на ось.
type Projection struct {
	Min, Max float64
}

// Main функция для проверки пересечения двух полигонов.
func PolygonsIntersect(p1, p2 Polygon) bool {
	// Получаем все оси для проверки
	axes := getAxes(p1)
	axes = append(axes, getAxes(p2)...)

	// Проверяем каждую ось
	for _, axis := range axes {
		// Проецируем оба полигона на ось
		proj1 := project(p1, axis)
		proj2 := project(p2, axis)

		// Проверяем, перекрываются ли проекции
		if !overlap(proj1, proj2) {
			return false
		}
	}

	return true
}

// getAxes возвращает все нормали к ребрам полигона.
func getAxes(p Polygon) []Point {
	axes := make([]Point, 0, len(p))
	for i := 0; i < len(p); i++ {
		p1 := p[i]
		p2 := p[(i+1)%len(p)] // Следующая точка (замыкаем полигон)

		// Вектор ребра
		edge := Point{p2.X - p1.X, p2.Y - p1.Y}

		// Нормаль к ребру (перпендикулярный вектор)
		normal := Point{-edge.Y, edge.X}
		axes = append(axes, normal)
	}
	return axes
}

// project проецирует полигон на ось и возвращает проекцию.
func project(p Polygon, axis Point) Projection {
	min := dot(p[0], axis)
	max := min

	for _, point := range p {
		proj := dot(point, axis)
		if proj < min {
			min = proj
		}
		if proj > max {
			max = proj
		}
	}

	return Projection{Min: min, Max: max}
}

// dot возвращает скалярное произведение двух точек.
func dot(p1, p2 Point) float64 {
	return p1.X*p2.X + p1.Y*p2.Y
}

// overlap проверяет, перекрываются ли две проекции.
func overlap(proj1, proj2 Projection) bool {
	return proj1.Min <= proj2.Max && proj2.Min <= proj1.Max
}

func main() {
	// Пример двух полигонов (выпуклых)
	polygon1 := Polygon{
		{0, 0},
		{4, 0},
		{4, 4},
		{0, 4},
	}
	polygon2 := Polygon{
		{2, 2},
		{6, 2},
		{6, 6},
		{2, 6},
	}

	// Проверяем пересечение
	if PolygonsIntersect(polygon1, polygon2) {
		fmt.Println("Полигоны пересекаются.")
	} else {
		fmt.Println("Полигоны не пересекаются.")
	}
}
