package main


import (
    "fmt"
    "strings"
    "strconv"
    //"io"
    "os"
    "bufio"
)

type vec2 struct {
    x, y float32
}

type vec3 struct {
    x, y, z float32
}

type vertex struct {
    pos vec3
    normal vec3
    tex vec2
}

type faceVtx struct {
    v, vt, vn int
}

func parseFloat(s string) float32 {
    v, err := strconv.ParseFloat(s, 32)
    if err != nil {
        panic(err)
    }
    return float32(v)
}

func parseInt(s string) int {
    v, err := strconv.ParseInt(s, 10, 0)
    if err != nil {
        panic(err)
    }
    return int(v)
}

func getInt(s []string, i int) int {
    if i >= len(s) || len(s[i]) == 0 {
        return 0
    }
    return parseInt(s[i])
}

func readWavefront(filename string) (vs []vec3, vts []vec2, vns []vec3, fs []faceVtx) {
    f, err := os.Open(filename)
    if err != nil { panic(err) }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    vs = make([]vec3, 0, 100)
    vns = make([]vec3, 0, 100)
    vts = make([]vec2, 0, 100)
    fs = make([]faceVtx, 0, 100)
    for scanner.Scan() {
        s := strings.Split(strings.Trim(scanner.Text(), " "), " ")
        if len(s) == 0 || s[0][0] == '#' {
            continue
        }
        switch s[0] {
            case "v":
                vs = append(vs, vec3{parseFloat(s[1]), parseFloat(s[2]), parseFloat(s[3])})
            case "vt":
                vts = append(vts, vec2{parseFloat(s[1]), parseFloat(s[2])})
            case "vn":
                vns = append(vns, vec3{parseFloat(s[1]), parseFloat(s[2]), parseFloat(s[3])})
            case "f":
                baseFaceIndex := len(fs)
                numVtxThisFace := 0
                for _, vtx := range s[1:] {
                    vtxIndices := strings.Split(vtx, "/")
                    v := faceVtx{ getInt(vtxIndices, 0), getInt(vtxIndices, 1), getInt(vtxIndices, 2)}
                    // Make a fan out of faces with more than 3 vertices
                    if numVtxThisFace >= 3 {
                        fs = append(fs, fs[baseFaceIndex])
                        fs = append(fs, fs[len(fs) - 2])
                    }
                    fs = append(fs, v)
                    numVtxThisFace++
                }
            default: ;
        }
    }
    return
}


func main() {
    if len(os.Args) < 2 {
        fmt.Println("Please supply the filename of a Wavefront OBJ file")
        os.Exit(1)
    }
    vs, vts, vns, fs := readWavefront(os.Args[1])
    fmt.Println("Read:", len(vs), "positions", len(vts), "texture coords", len(vns), "normals", len(fs)/3, "faces")
    fmt.Println("capacities:", cap(vs), "positions", cap(vts), "texture coords", cap(vns), "normals", cap(fs), "faces")

    // We now have read all the mesh data, but each face vertex has different indices for pos, normal and texture
    // we need to coalesce complete vertices

    // Use a map to track which complete vertices we have already added to the vertex buffer
    // The key is the vertex, the value is the index of that vertex in the vertex buffer
    vertexMap := make(map[vertex]int, len(fs))
    vertices := make([]vertex, 0, len(fs)) // vertex buffer
    indices := make([]int, 0, len(fs)) // index buffer
    for _, fv := range fs {
        vtx := new(vertex)
        // Indices can be negative, and positive ones are 1-based
        if fv.v != 0 { if fv.v > 0 { vtx.pos = vs[fv.v-1] } else { vtx.pos = vs[len(vs)+fv.v] } }
        if fv.vt != 0 { if fv.vt > 0 { vtx.tex = vts[fv.vt-1] } else { vtx.tex = vts[len(vts) + fv.vt] } }
        if fv.vn != 0 { if fv.vn > 0 { vtx.normal = vns[fv.vn-1] } else { vtx.normal = vns[len(vns) + fv.vn] } }

        existingIndex, ok := vertexMap[*vtx]
        if !ok {
            existingIndex := len(vertices)
            vertices = append(vertices, *vtx)
            vertexMap[*vtx] = existingIndex
        }
        indices = append(indices, existingIndex)
    }
    fmt.Println("Coalesced:", len(vertices), "vertices", len(indices)/3, "faces")
}
