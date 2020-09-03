package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

type Vertex struct {
	vert int
}

type Edge struct {
	from, to Vertex
}

type Graph struct {
	vertices         map[Vertex]bool //Vertices
	edges            map[Edge]bool   //Edges
	numvert, numedge int             //Number of Vertices/Edges
	directed         bool            //true if graph is directed
}

//a type for the cfi gadgets
type Cfigadget struct {
	inner      []Vertex //inner vertices
	outer      []Vertex //outer vertices
	ref        *Vertex  //referenced vertex of the underlying graph
	degree     int      //degree of the referenced vertex
	edges      []Edge   //edges inside the gadget
	neighbours []Vertex //neighbours of ref in the underlying graph
}

//add a new vertex to the graph
func (g *Graph) addVertex() (Vertex, error) {
	v := Vertex{
		vert: g.numvert + 1,
	}
	g.vertices[v] = true
	g.numvert++
	return v, nil
}

//adds a number of vertices to the graph
func (g *Graph) addVertices(vertices int) ([]Vertex, error) {
	i := 1
	var arr []Vertex
	arr = make([]Vertex, vertices)

	for i <= vertices {
		v, err := g.addVertex()
		if err != nil {
			log.Fatal(err)
		}
		arr[i-1] = v
		i++
	}

	return nil, nil
}

//add an edge from i to j
func (g *Graph) addEdge(i, j Vertex) error {
	e := Edge{
		from: i,
		to:   j,
	}
	eback := Edge{
		from: j,
		to:   i,
	}

	//If the Edge already exists in one of the directions then don't add it
	if g.edges[e] || g.edges[eback] {
		fmt.Println("Kante von ", e.from.vert, " nach ", e.to.vert, " bereits vorhanden")
		return nil
	}
	//Also if the vertices don't even exist, don't add the edge and throw an error
	if !(g.vertices[i] && g.vertices[j]) {
		//log.Fatal("Knoten existieren nicht, Kante wird nicht hinzugefügt")
	}

	g.edges[e] = true

	g.numedge++
	return nil
}

func (g *Graph) removeEdge(e Edge) error {
	g.edges[e] = false
	return nil
}

//getter for the vertices
func (g *Graph) getVertices() ([]Vertex, error) {
	var arr []Vertex
	arr = make([]Vertex, g.numvert)

	i := 0
	for node, value := range g.vertices {
		if value {
			arr[i] = node
			i++
		}
	}

	return arr, nil
}

func (g *Graph) degree(v Vertex) (int, error) {
	deg := 0

	//	if g.vertices[v] {
	//		fmt.Println("Knoten existiert nicht im Graphen")
	//		return -1, nil
	//	}
	for edge, value := range g.edges {
		if (edge.from.vert == v.vert || edge.to.vert == v.vert) && value {
			deg++
		}
	}

	return deg, nil
}

//returns the neighbourhood of v in some graph g
func (g *Graph) neighbourhood(v Vertex) ([]Vertex, error) {
	deg, err := g.degree(v)
	if err != nil {
		log.Fatal(err)
	}

	var neighbours []Vertex
	neighbours = make([]Vertex, deg)

	i := 0
	for edge, value := range g.edges {
		if (edge.from.vert == v.vert) && value {
			neighbours[i] = edge.to
			i++
		} else if (edge.to.vert == v.vert) && value {
			neighbours[i] = edge.from
			i++
		}
	}

	return neighbours, nil
}

func buildGraph(dir bool) (Graph, error) {

	g := Graph{
		vertices: make(map[Vertex]bool),
		edges:    make(map[Edge]bool),
		directed: dir,
		numvert:  0,
		numedge:  0,
	}

	return g, nil
}

//Write graph into output
func (g *Graph) writeGraph(output io.Writer) error {
	header := "p tw " + strconv.Itoa(g.numvert) + " " + strconv.Itoa(g.numedge) + "\n"
	io.WriteString(output, header)

	line := ""
	for edge, value := range g.edges {
		if value {
			line = strconv.Itoa(edge.from.vert) + " " + strconv.Itoa(edge.to.vert) + "\n"
			io.WriteString(output, line)
		}
	}

	return nil
}

//Print graph
func (g *Graph) printGraph() error {
	fmt.Println("p tw ", g.numvert, " ", g.numedge)

	for edge, value := range g.edges {
		if value {
			fmt.Println(edge.from.vert, " ", edge.to.vert)
		}
	}

	return nil
}

func (g *Cfigadget) printGadget() error {
	//fmt.Println("p tw ", g.numvert, " ", g.numedge)

	for _, edge := range g.edges {
		fmt.Println(edge.from.vert, " ", edge.to.vert)
	}

	return nil
}

//returns 1+2+3+...+i
func sum(i int) int {
	n := 0
	for i > 0 {
		n += i
		i--
	}
	return n
}

//build the grid graph of size ixi and write into output
func makegrid(i int, g *Graph) error {

	nodes := i * i         //number of vertices
	edges := i * (2*i - 2) //number of edges

	if g.numvert < nodes {
		_, err := g.addVertices(nodes - g.numvert)
		if err != nil {
			log.Fatal(err)
		}
	}
	if g.numvert > nodes {
		log.Fatal("Too many vertices")
	}

	from := 1
	to := 2
	run := 0
	for edges > 0 {
		a := Vertex{vert: from}
		b := Vertex{vert: to}
		g.addEdge(a, b)
		edges--
		from++
		to++
		if from%i == 0 && run == 0 {
			from++
			to++
		}
		if to > nodes {
			from = 1
			to = 1 + i
			run = 1
		}

	}

	return nil
}

//build rectangular ixj grid
func makerectangle(i int, j int, g *Graph) error {
	nodes := i * j //number of vertices

	if g.numvert < nodes {
		_, err := g.addVertices(nodes - g.numvert)
		if err != nil {
			log.Fatal(err)
		}
	}
	if g.numvert > nodes {
		log.Fatal("Too many vertices")
	}

	from := 1
	to := 2
	run := false
	done := false
	for !done {
		a := Vertex{vert: from}
		b := Vertex{vert: to}
		g.addEdge(a, b)
		from++
		to++
		if from%j == 0 && !run {
			from++
			to++
		}
		if run && to > nodes {
			done = true
		}
		if to > nodes {
			from = 1
			to = 1 + j
			run = true
		}
	}
	return nil
}

//build the wall graph
func makewall(i int, g *Graph) error {
	nodes := 2 * i * i
	edges := i*(2*i-2) + i*i

	if g.numvert < nodes {
		_, err := g.addVertices(nodes - g.numvert)
		if err != nil {
			log.Fatal(err)
		}
	}
	if g.numvert > nodes {
		log.Fatal("Too many vertices")
	}

	from := 1
	to := 2
	for edges > 0 {
		a := Vertex{vert: from}
		b := Vertex{vert: to}
		g.addEdge(a, b)
		edges--

		//add the vertical edges in the wall graph
		//only if the index isn't outside of the graph
		lower := (from + 2*i + 1)
		if from%2 == 1 && lower <= nodes {
			b = Vertex{vert: lower}
			g.addEdge(a, b)
			edges--
		}

		from++
		to++
		if from%(2*i) == 0 {
			from++
			to++
		}

	}

	return nil
}

//build a hexagonal grid of radius i
func makehex(radius int, g *Graph) error {
	nodes := radius * radius * 6 //number of vertices

	if g.numvert < nodes {
		_, err := g.addVertices(nodes - g.numvert)
		if err != nil {
			log.Fatal(err)
		}
	}
	if g.numvert > nodes {
		log.Fatal("Too many vertices")
	}

	i := 1
	//vertices, err := g.getVertices()
	//if err != nil {
	//	log.Fatal(err)
	//}

	for i <= radius {
		//draw a circle for each radius
		currad := i * i * 6                      //circumference of the current circle
		startcirc := ((i - 1) * (i - 1) * 6) + 1 //start value of the current circle
		j := startcirc

		for j < currad {
			a := Vertex{vert: j}
			b := Vertex{vert: j + 1}
			g.addEdge(a, b)
			fmt.Println(j)
			j++
		}
		a := Vertex{vert: currad}
		b := Vertex{vert: startcirc}
		g.addEdge(a, b)

		//Connect the current circle with the one before
		if i > 1 {
			innerpos := (i-2)*(i-2)*6 + 1           //Start value of the inner circle
			innercirc := (i - 1) * (i - 1) * 6      //circumference of the inner circle
			outerpos := ((i - 1) * (i - 1) * 6) + 1 //start value of the outer circle

			for innerpos <= innercirc {
				//if innerpos%2 == 0 {

				v := Vertex{vert: innerpos}
				deg, err := g.degree(v)
				if err != nil {
					log.Fatal(err)
				}

				/*In a hexagonal grid, all inner vertices have degree 3,
				so if a vertex of the inner circle does not have this
				degree, it lacks connection to the outer circle*/
				if deg != 3 {
					a := Vertex{vert: innerpos}
					b := Vertex{vert: 4 + 3*innerpos}
					/*Funktioniert so bisher nur bei Kreisen mit Radius 1 und 2*/
					g.addEdge(a, b)
					outerpos += 3
				}
				innerpos++
			}
		}

		i++
	}

	return nil
}

//build a pyramid grid of base width i
func makepyramid(basis int, g *Graph) error {

	nodes := sum(basis)
	edges := 3 * sum(basis-1)

	if g.numvert < nodes {
		_, err := g.addVertices(nodes - g.numvert)
		if err != nil {
			log.Fatal(err)
		}
	}
	if g.numvert > nodes {
		log.Fatal("Too many vertices")
	}

	vone := 0
	vtwo := 1
	vthree := 2
	level := 1
	for edges > 0 {
		//add triangles for one level
		for i := level; i > 0; i-- {
			vone++
			vtwo++
			vthree++

			//fmt.Println(vone, " ", vtwo, " ", vthree, " Level: ", level)

			a := Vertex{vert: vone}
			b := Vertex{vert: vtwo}
			c := Vertex{vert: vthree}
			g.addEdge(a, b)
			g.addEdge(a, c)
			g.addEdge(b, c)
			edges -= 3
		}

		vone = sum(level)
		level++
		vtwo = vone + level
		vthree = vtwo + 1
	}

	return nil
}

//build the expanded grids G^+
func expandgraph(g *Graph) error {

	newedges := make(map[Edge]bool)

	for e := range g.edges {
		from := e.from
		to := e.to

		v1, err := g.addVertex()
		if err != nil {
			log.Fatal(err)
		}
		v2, err := g.addVertex()
		if err != nil {
			log.Fatal(err)
		}

		e1 := Edge{
			from: from,
			to:   v1,
		}
		e2 := Edge{
			from: v1,
			to:   v2,
		}
		e3 := Edge{
			from: v2,
			to:   to,
		}

		newedges[e1] = true
		newedges[e2] = true
		newedges[e3] = true

		err = g.removeEdge(e)
		if err != nil {
			log.Fatal(err)
		}

	}

	g.edges = newedges

	return nil
}

//apply the cfi construction to a graph
func cfi(g *Graph) (Graph, error) {
	cfig, err := buildGraph(false)
	if err != nil {
		log.Fatal(err)
	}

	//**make an array with all the Cfi Gadgets**//
	var gadgets map[Vertex]Cfigadget
	gadgets = make(map[Vertex]Cfigadget)

	//number the vertices in the original graphs
	vertcount := 1

	//the current vertcount will be the key for a vertex of the new graph
	var nodekeys map[int]Vertex
	nodekeys = make(map[int]Vertex)

	for node, value := range g.vertices {
		if value {
			deg, err := g.degree(node)
			if err != nil {
				log.Fatal(err)
			}

			//Make a new Gadget
			gadget := Cfigadget{
				outer:      make([]Vertex, 2*deg),
				inner:      make([]Vertex, int(math.Exp2(float64(deg-1)))),
				ref:        &node,
				degree:     deg,
				edges:      make([]Edge, deg*int(math.Exp2(float64(deg-1)))),
				neighbours: make([]Vertex, deg),
			}

			gadget.neighbours, err = g.neighbourhood(node)

			//insert actual vertices into the gadget
			for ind := range gadget.inner {
				gadget.inner[ind] = Vertex{vert: vertcount}
				nodekeys[vertcount] = gadget.inner[ind]
				vertcount++
			}
			for ind := range gadget.outer {
				gadget.outer[ind] = Vertex{vert: vertcount}
				nodekeys[vertcount] = gadget.outer[ind]
				vertcount++
			}

			//just somewhere to save the number of edges to
			edgecount := 0

			//Build the gadget
			max := int(math.Exp2(float64(deg)))
			innercount := 0 //tells how many of the inner vertices have been used
			for i := 0; i < max; i++ {
				cur := i

				//save the calculated bits
				var bits []int
				bits = make([]int, deg)

				//number of 1s in the binary representation of the number
				ones := 0

				//calculate binary representation
				for j := 0; j < deg; j++ {
					if cur&1 == 1 {
						ones++
						bits[j] = 1
					} else {
						bits[j] = 0
					}
					cur = cur >> 1
				}

				//with an even number of ones, connect the vertices
				if ones%2 == 0 {
					for ind := 0; ind < deg; ind++ {
						outbit := ((ind + 1) * 2) / 2 //which of the outer bits is controlled

						if bits[outbit-1] == 0 { //connect to b vertices
							outvertex := ind*2 + 1
							gadget.edges[edgecount] = Edge{
								from: gadget.outer[outvertex],
								to:   gadget.inner[innercount],
							}
						} else { //connect to a vertices
							outvertex := ind * 2
							gadget.edges[edgecount] = Edge{
								from: gadget.outer[outvertex],
								to:   gadget.inner[innercount],
							}
						}
						edgecount++
					}
					innercount++
				}

			}

			//Save it to the map
			gadgets[node] = gadget
		}
	}

	//**build a graph from the gadgets**//

	//add inner and outer vertices (the content of nodekeys) to the graph
	for _, v := range nodekeys {
		cfig.vertices[v] = true
		cfig.numvert++
	}

	//add the known edges from within the gadgets
	for _, gadget := range gadgets {
		for _, e := range gadget.edges {
			cfig.addEdge(e.from, e.to)
		}
	}

	//connect the gadgets
	for e := range g.edges {
		edgefrom := e.from
		edgeto := e.to
		gadgetfrom := gadgets[edgefrom]
		gadgetto := gadgets[edgeto]

		neighbournofrom := 0 //number of the to-neighbour in the "from" gadget
		neighbournoto := 0   //number of the from-neighbour in the "to" gadget

		for ind, n := range gadgetfrom.neighbours {
			if n == edgeto {
				neighbournofrom = ind
				break
			}
		}
		for ind, n := range gadgetto.neighbours {
			if n == edgefrom {
				neighbournoto = ind
				break
			}
		}
		cfig.addEdge(gadgetfrom.outer[2*neighbournofrom], gadgetto.outer[2*neighbournoto])
		cfig.addEdge(gadgetfrom.outer[2*neighbournofrom+1], gadgetto.outer[2*neighbournoto+1])
	}

	//cfig.printGraph()
	return cfig, nil
}

//test some construction and print the graph to console
func testgraphs() {
	g, err := buildGraph(false)
	if err != nil {
		log.Fatal(err)
	}

	err = makegrid(2, &g)
	if err != nil {
		log.Fatal(err)
	}

	err = expandgraph(&g)
	if err != nil {
		log.Fatal(err)
	}

	err = g.printGraph()
	if err != nil {
		log.Fatal(err)
	}
}

//build graphs and write the files to output
func writegraphs() {

	for i := 2; i <= 10; i++ {
		g, err := buildGraph(false)
		if err != nil {
			log.Fatal(err)
		}
		err = makegrid(i, &g)
		if err != nil {
			log.Fatal(err)
		}
		//g, err = cfi(&g)
		//if err != nil {
		//	log.Fatal(err)
		//}

		expandgraph(&g)

		filename := "expgrid" + strconv.Itoa(i) + ".gr"
		fmt.Println(filename, " created")
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		err = g.writeGraph(file)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	writegraphs()
	//testgraphs()
}
