package main

import (
	"fmt"
	"time"
	"log"
	"net/http"
	// "encoding/json"
	"math"
	"sort"
	// "errors"

	"github.com/gorilla/websocket"
	"encoding/binary"
)

type circleT struct{
	xx		uint16
	yy		uint16
	rayon	uint16
}

type snkeT struct{
	id    uint16
	lvAmt uint16
	xx    uint16
	yy    uint16
	lnpxx uint16
	lnpyy uint16
	angle float32
	sin		float64
	cos		float64
	ehang float32
	width	uint16
	sp    uint16	//should be speed
	score uint16
}

type snkHeadT struct{
	id    		uint16
	xx    		uint16
	yy    		uint16
	sz    		uint16
	dist  		float64
}

type snkBodyT struct{
	id    		uint16
	xx    		uint16
	yy    		uint16
	sz    		uint16
	dist  		float64
}

type foodT struct{
	xx    		uint16
	yy    		uint16
	sz    		uint16
	dist  		float64
	angle 		float64
	secto		uint8
	tooClose	bool
}

type bySize []foodT

func (s bySize) Len() int {
	return len(s)
}
func (s bySize) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s bySize) Less(i, j int) bool {
	return s[i].dist < s[j].dist
}

// type sectoT struct{
// 	foodT,
// 	sum
// }

func getIdx(dx,dy int32) uint8 {
  if dy>=0 {                    // 1/2
    if dx>=0 {                  // 1/4
      if dx>=dy {               // 1/8
        dy2:=int32(float32(dy)*2.414214)
        if dx>dy2 {    					// 1/16
					return 0
				} else {    						// 2/16
					return 1
				}
      }else{                    // 2/8
        dx2:=int32(float32(dx)*2.414214)
        if dy>dx2 {    					// 4/16
					return 3
				}else{    							// 3/16
					return 2
				}
      } 
    }else{                      // 2/4
      dx=-dx
      if dy>=dx {               // 3/8
        dx2:=int32(float32(dx)*2.414214)
        if dy>dx2 {    					// 5/16
					return 4
				}else{    							// 6/16
				  return 5
				}
      }else{                    // 4/8
        dy2:=int32(float32(dy)*2.414214)
        if dx>dy2 {    					// 8/16
					return 7
				}else{    							// 7/16
				  return 6
				}
      }
    }
  }else{                        // 2/2
    dy=-dy
    if dx>=0 {                  // 4/4
      if dx>=dy {               // 8/8
        dy2:=int32(float32(dy)*2.414214)
        if dx>dy2 {   					// 16/16
					return 15 
				}else{   								// 15/16
					return 14
				}
      }else{          					// 7/8
        dx2:=int32(float32(dx)*2.414214)
        if dy>dx2 {
					return 12   					// 13/16
				}else{
					return 13   					// 14/16
				}
      }
    }else{                      // 3/4
      dx=-dx
      if dy>=dx {               // 6/8
        dx2:=int32(float32(dx)*2.414214)
        if dy>dx2 {   					// 9/16
					return 11
				}else{   								// 10/16
				  return 10
				}
      }else{                    // 5/8
        dy2:=int32(float32(dy)*2.414214)
        if dx>dy2 {   					// 12/16
					return  8
				}else{   								// 11/16
					return  9
				}
      }
    }
  }
}

func inCercle(circle circleT,xx,yy uint16) bool {
	if xx>circle.xx+circle.rayon {return false}
	if xx<circle.xx-circle.rayon {return false}
	if yy>circle.yy+circle.rayon {return false}
	if yy<circle.yy-circle.rayon {return false}
	var dx int32 =int32(circle.xx)-int32(xx)
	var dy int32 =int32(circle.yy)-int32(yy)
	if dx*dx+dy*dy<=int32(circle.rayon)*int32(circle.rayon) {
		// fmt.Println("On elimine")
		return true}
	return false
}

//takes at entering snkes array, populating snkHeads slice
//except bodies points, no other snk informations will be considered later
func getSnkHeads(snkes []snkeT,snkHeads []snkHeadT,nb int,mySnk snkeT) bool{
	var sx int32 =int32(mySnk.xx)
	var sy int32 =int32(mySnk.yy)
	for j:=0;j<16;j++{
		snkHeads[j].sz=0
	}
	// Here we would later analyse the direction of other snk
	for i:=0;i<nb;i++{
		var s=snkes[i]
		if s.id==mySnk.id{continue}
		if s.lvAmt!=1{continue}
		var dx int32 =sx-int32(s.xx)
		var dy int32 =sy-int32(s.yy)
		dist:=math.Sqrt(float64(dx*dx+dy*dy))
		idx:=getIdx(dx,dy)
		if snkHeads[idx].sz==0 {
			snkHeads[idx].id=s.id
			snkHeads[idx].xx=s.xx
			snkHeads[idx].yy=s.yy
			snkHeads[idx].sz=1
			snkHeads[idx].dist=dist
		}else{
			if dist<snkHeads[idx].dist {
				snkHeads[idx].id=s.id
				snkHeads[idx].xx=s.xx
				snkHeads[idx].yy=s.yy
				// snkHeads[idx].sz=1
				snkHeads[idx].dist=dist
			}
		}
	}
	return false
}

func getSnkBodys(snksPts [][]snkBodyT,secPts [16][]snkBodyT,nb int,mySnk snkeT) bool{
	// var sx int32 =int32(mySnk.xx)
	// var sy int32 =int32(mySnk.yy)
	for j:=0;j<16;j++{
		secPts[j]=nil
	}
	for i:=0;i<nb;i++{
		var snk=snksPts[i]
		// for j:=0;j<nb;j++{

		// }
	}
	return false
}

func checkCollision(nb int) bool{

	return false
}

func prepareFood(snake snkeT,foods []foodT,nbFoods int,log bool) foodT {
	type sectoT struct{
		size 		float64
		foodId	uint16
		dist		float64
	}
	var sectoSiz [16] sectoT
	var i uint16 =0
	var maxIdx uint8 =0
	var maxSz float64 =0
	var sideCircleR,sideCircleL circleT

	// On calcule les cercles d'exclusion latéraux
	var snSin=uint16(snake.sin*float64(snake.width))
	var snCos=uint16(snake.cos*float64(snake.width))
	sideCircleR.xx=snake.lnpxx-snSin
	sideCircleR.yy=snake.lnpyy+snCos
	sideCircleR.rayon=snake.width
	sideCircleL.xx=snake.lnpxx+snSin
	sideCircleL.yy=snake.lnpyy-snCos
	sideCircleL.rayon=snake.width

	// var init bool =true

	for i:=0;i<nbFoods;i++ {
		var fd foodT =foods[i]
		if inCercle(sideCircleL,fd.xx,fd.yy) {foods[i].tooClose=true;continue}
		if inCercle(sideCircleR,fd.xx,fd.yy) {foods[i].tooClose=true;continue}
		foods[i].tooClose=false
		var dx int32 =int32(fd.xx)-int32(snake.xx)
		var dy int32 =int32(fd.yy)-int32(snake.yy)
		// foods[i].angle=math.Atan2(float64(dy),float64(dx))
		foods[i].dist=math.Sqrt(float64(dx*dx+dy*dy))
		idx:=getIdx(dx,dy)
		foods[i].secto=idx
	}
	sort.Sort(bySize(foods))	// tri en distances croissantes

	for i=0;i<16;i++ {
		sectoSiz[i].size=0
		sectoSiz[i].dist=0
	}
	for i:=0;i<nbFoods;i++ {
		var food foodT =foods[i]
		idx:=food.secto
		if food.tooClose {continue}

		// var ddist float64 =foods[i].dist-sectoSiz[idx].dist
		sectoSiz[idx].dist=food.dist // a quoi ca sert ?
		// sectoSiz[idx].size+=float64(foods[i].sz)/ddist
		// sectoSiz[idx].size+=float64(foods[i].sz*foods[i].sz)/(ddist*foods[i].dist)

		// sectoSiz[idx].size+=float64(food.sz)/food.dist// acceptable mais change trop
		sectoSiz[idx].size+=float64(food.sz)/food.dist/food.dist // mieux, plus stable
									// mais ne va pas chercher les grosses valeurs

	}
	// if log {
	// 	fmt.Printf("snkXX=%05d snkYY=%05d\n",snake.xx,snake.yy)

	// 	for i:=0;i<10;i++ {
	// 		var fd foodT =foods[i]
	// 		fmt.Printf("X=%05d Y=%05d snX=%05d snY=%05d dist=%05d\n",fd.xx,fd.yy,snake.xx,snake.yy,fd.dist)
	// 	}
	// }
	// if log {
		// fmt.Println()
		// for i:=0;i<5;i++ {
		// 	var fd foodT =foods[i]
		// 	fmt.Printf("X=%05d Y=%05d secto=%05d sz=%05d dist=%05d\n",fd.xx,fd.yy,fd.secto,fd.sz,fd.dist)
		// }
	// }

	// maxIdx=0
	maxSz=0
	for i=0;i<16;i++ {
		// fmt.Printf("[%01d] %05f sz=%05f\n",i,sectoSiz[i],maxSz)
		if sectoSiz[i].size>maxSz {
			maxSz=sectoSiz[i].size;maxIdx=uint8(i)
		}
	}
	// fmt.Printf("\nbest secto %05d sz=%05f\n",maxIdx,maxSz)

	var idx int
	for idx=0;idx<nbFoods;idx++ {
		var fd foodT =foods[idx]
		if fd.tooClose {continue}
		// if fd.dist>100*100{break}
		// if fd.secto==maxIdx&&fd.dist>100*100{break}
		if fd.secto==maxIdx{break}
	}
	if(idx==nbFoods){
		fmt.Printf("\n-------------Erreur %05d %05d\n",nbFoods,idx)
		fmt.Printf("\nbest secto %05d sz=%05f\n",maxIdx,maxSz)
		for i:=0;i<10 && i<nbFoods;i++ {
			var fd foodT =foods[i]
			fmt.Printf("X=%05d Y=%05d secto=%05d sz=%05d dist=%05f tooClose=%t\n",fd.xx,fd.yy,fd.secto,fd.sz,fd.dist,fd.tooClose)
		}
		for i=0;i<16;i++ {
			fmt.Printf("[%3d] sz=%05f dist=%5.5f\n",i,sectoSiz[i].size,sectoSiz[i].dist)
		}
	}

	// fmt.Printf("prep: fX=%05d fY=%05d sX=%05d Sy=%05d dist=%05d\n",foods[idx].xx,foods[idx].yy,snake.xx,snake.yy,foods[idx].dist)
	return foods[idx]
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	var timeStart,timeEnd time.Time
	// var timeEnd time.Time

	var nbFoods int =0
	var nbSnkes int =0

	var stats uint16 =0
	var foods []foodT
	var snkes []snkeT
	var snksPts [][]snkBodyT
	var snkIdx uint16

	// var dist  uint16 =0
	var i int =0
	var comd uint16
	var init bool =true
	var snake snkeT
	var circleAngle float64 =0
	var circleAngPas float64 =0.1
	var strategie uint8 =1
	var tpsTrame int
	var tpsEnvoi int

	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// print out that message for clarity
		// log.Println(string(p))
		// fmt.Printf("Message type %T \n",p)
		// fmt.Printf("Message type %T \n",p[0])
	comd=binary.LittleEndian.Uint16(p)
		// if init&&comd==102 {
		// 	var toto1 uint16
		// 	var toto2 uint16
		// 	var toto3 uint16
		// 	for i=0;i<10;i++ {
		// 		toto1=binary.LittleEndian.Uint16(p[6*i+2:6*i+4])
		// 		toto2=binary.LittleEndian.Uint16(p[6*i+4:6*i+6])
		// 		toto3=binary.LittleEndian.Uint16(p[6*i+6:6*i+8])
		// 		fmt.Printf("i=%05d X=%05d Y=%05d Size=%05d\n",i,toto1,toto2,toto3)
		// 	}
		// 	init=false
		// }
	
	switch(comd){
		case 101:
			timeStart=time.Now()
			snake.xx=binary.LittleEndian.Uint16(p[2:4])
			snake.yy=binary.LittleEndian.Uint16(p[4:6])
			snake.lnpxx=binary.LittleEndian.Uint16(p[6:8])
			snake.lnpyy=binary.LittleEndian.Uint16(p[8:10])
			snake.width=binary.LittleEndian.Uint16(p[10:12])			
			var angIn=binary.LittleEndian.Uint16(p[12:14])
			if angIn>32768 {
				angIn=^angIn+1
				snake.angle=-float32(angIn)/20805.0
			}else{
				snake.angle=float32(angIn)/20805.0
			}
			snake.sin,snake.cos=math.Sincos(float64(snake.angle))
			snake.id=binary.LittleEndian.Uint16(p[14:16])
			tpsTrame=(int)(binary.LittleEndian.Uint16(p[16:18]))
			tpsEnvoi=(int)(binary.LittleEndian.Uint16(p[18:20]))
			fmt.Printf("ms trame:%5d Envoi:=%5d\n",tpsTrame,tpsEnvoi)

			// if !init {
				// fmt.Printf("snX=%05d snY=%05d NbFood=%05d Stats=%05d\r",	snake.xx,	snake.yy,nbFoods,stats)
				// fmt.Printf("\n 101 snX=%05d snY=%05d width=%05d angle=%05f\n",snake.xx,snake.yy,snake.width,snake.angle)
			// }

		case 103:
			if strategie==1 {
				strategie=2
				log.Println("\nCircle begin")
			} else {
				strategie=1
				log.Println("\nCircle end")
			}

		case 104:
			if strategie==2 {
				log.Println("\nCircle increase %f",circleAngPas)
				circleAngPas/=1.1
			}

		case 105:
			if strategie==2 {
				if circleAngPas<=0.1950{
					log.Println("\nCircle decrease %f",circleAngPas)
					circleAngPas*=1.1
				}else{
					log.Println("\nCircle decrease is min")
				}
			}

		case 102:		// arrivée foods
			// fmt.Printf("Inside comd102 %d\n",nbFoods)
			nbFoods=(int)(binary.LittleEndian.Uint16(p[2:4]))

			foods = make([]foodT,nbFoods) 	// seulement ici pour ne pas reset les anciens
			for i=0;i<nbFoods;i++ {
				foods[i].xx=binary.LittleEndian.Uint16(p[6*i+4:6*i+ 6])
				foods[i].yy=binary.LittleEndian.Uint16(p[6*i+6:6*i+ 8])
				foods[i].sz=binary.LittleEndian.Uint16(p[6*i+8:6*i+10])
			}

			stats=0
			for i=0;i<nbFoods;i++ {
				var getstat uint16 =0
				// getstat=binary.LittleEndian.Uint16(p[6*i+6:6*i+8])
				getstat=foods[i].sz
				stats+=getstat
				// fmt.Printf("Status i=%05d getstat=%05d Stats=%05d\r",i,getstat,stats)
			}

		case 106:			// arrivée liste serpents
			nbSnkes=(int)(binary.LittleEndian.Uint16(p[2:4]))
			// fmt.Printf("nbSnkes %d\n",nbSnkes)
			snkes = make([]snkeT,nbSnkes) // why not only an array?
			snksPts = make([][]snkBodyT,nbSnkes)
			snkIdx=0;

			// snkes = [nbSnkes]snkeT
			for i=0;i<nbSnkes;i++ {
				snkes[i].id			=binary.LittleEndian.Uint16(p[10*i+ 4:10*i+ 6])
				snkes[i].lvAmt	=binary.LittleEndian.Uint16(p[10*i+ 6:10*i+ 8])
				snkes[i].xx			=binary.LittleEndian.Uint16(p[10*i+ 8:10*i+10])
				snkes[i].yy			=binary.LittleEndian.Uint16(p[10*i+10:10*i+12])
				// snkes[i].ehang	=binary.LittleEndian.Uint16(p[10*i+12:10*i+14])
			}

		case 107:			// Arrivée points premier serpent
			var id uint16
			var nb uint16
			id	=binary.LittleEndian.Uint16(p[2:4])
			nb	=binary.LittleEndian.Uint16(p[4:6])
			fmt.Printf("107  id:%5d nb:=%5d\n",id,nb)
			snksPts[0]=make([]snkBodyT,nb);
			snkIdx=1;

		case 108:			// Arrivée points autres serpents
			var id uint16
			var nb uint16
			id	=binary.LittleEndian.Uint16(p[2:4])
			nb	=binary.LittleEndian.Uint16(p[4:6])
			fmt.Printf("108  id:%5d nb:=%5d\n",id,nb)
			snksPts[snkIdx]=make([]snkBodyT,nb);
			snkIdx++;

		case 109:			// Arrivée points dernier serpent
			var id uint16
			var nb uint16
			id	=binary.LittleEndian.Uint16(p[2:4])
			nb	=binary.LittleEndian.Uint16(p[4:6])

			// if snkIdx!=uint16(nbSnkes) {
			// 	panic(errors.New("Last body is not last body"))
			// }

			fmt.Printf("109  id:%5d nb:=%5d\n\n",id,nb)
			snksPts[snkIdx]=make([]snkBodyT,nb);
			snkIdx++;
		}

		// Traitement et envoi
		if comd==102{
			timeEnd=time.Now()
			diff := timeEnd.Sub(timeStart)
			fmt.Printf("Microseconds: %d\n", diff.Nanoseconds()/1000)

			// var toSend[2] uint16
			toSend:=make([]byte,4)
			if strategie==2 {	//cercle
				// toSend[0]=uint16(100*math.Sin(circleAngle))+snake.xx
				var val1 uint16=uint16(100*math.Sin(circleAngle))+snake.xx
				// toSend[1]=uint16(100*math.Cos(circleAngle))+snake.yy
				var val2 uint16=uint16(100*math.Cos(circleAngle))+snake.yy
				binary.LittleEndian.PutUint16(toSend,val1)
				binary.LittleEndian.PutUint16(toSend[2:4],val2)

				circleAngle+=circleAngPas
				if circleAngle>math.Pi{circleAngle-=2*math.Pi}
			}else{ //action
				var goFood foodT
				var secPts [16][]snkBodyT

				// Pre-analyze snakes
				snkHeads := make([]snkHeadT, 16)
				getSnkHeads(snkes,snkHeads,nbSnkes,snake)
				getSnkBodys(snkes,secPts,nbSnkes,snake)

				// Eating
				if nbFoods>0 {
					goFood=prepareFood(snake,foods,nbFoods,init)
					// toSend[0]=goFood.xx
					// toSend[1]=goFood.yy
					binary.LittleEndian.PutUint16(toSend,goFood.xx)
					binary.LittleEndian.PutUint16(toSend[2:4],goFood.yy)
				}
			}

			init=false
			// dist=goFood.dist
			// fmt.Printf("send: X=%05d Y=%05d Size=%05d dist=%05d\n",goFood.xx,goFood.yy,goFood.sz,goFood.dist)

			// onJson, _ := json.Marshal(toSend)
			// err = conn.WriteJSON(toSend)
			// if err != nil {
			// 	log.Println("WriteMessage:", err)
			// 	return
			// }
				// code fonctionnel pour version binaire directe
			// toSend:=make([]byte,4)
			// binary.LittleEndian.PutUint16(toSend, 10106)
			// binary.LittleEndian.PutUint16(toSend[2:4], 10107)
			// fmt.Println(toSend)
			// err := ws.WriteMessage(websocket.byte, []byte(t.String()))
			err = conn.WriteMessage(websocket.BinaryMessage,toSend)
			if err != nil {
				log.Println("write:", err)
				return
			}
		}

	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page for Slither Summary")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	// err = ws.WriteMessage(1, []byte("Hi Client!"))
	// if err != nil {
	// 	log.Println(err)
	// }

	// code fonctionnel pour version binaire directe
	toSend:=make([]byte,4)
	binary.LittleEndian.PutUint16(toSend, 10106)
	binary.LittleEndian.PutUint16(toSend[2:4], 10107)
	fmt.Println(toSend)
	// err := ws.WriteMessage(websocket.byte, []byte(t.String()))
	err = ws.WriteMessage(websocket.BinaryMessage,toSend)
	if err != nil {
		log.Println("write:", err)
		return
	}

	//version JSON fonctionnelle
	// var toSend[2] uint16
	// toSend[0]=10106
	// toSend[1]=10107
	// // onJson, _ := json.Marshal(toSend)
	// err = ws.WriteJSON(toSend)
	// if err != nil {
	// 	log.Println("WriteMessage:", err)
	// 	return
	// }

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Starting server")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
