package main

//Created By: Kyle McClendon

import (
	"fmt" //Printing
	"os" //Reading
	"log" //Logging Problems
	"bufio" //Reading Strings
	"io" //EOF Stuff
	"strings" //String Handling
	"strconv" //Convert Strings To Decimals
	"unicode/utf8" //Convert Characters To Decimal Values
	)

//Converts Characters To Decimal Values ("A" = 65)
func stringConvert(dlls string) []int{
	dllRefs := make([]int, len(dlls));
	chars := strings.Split(dlls, "")
	for i := 0; i < len(chars); i++ {
		r, size := utf8.DecodeRuneInString(chars[i]);
		if(size != 0){
			dllRefs[i] = (int(r) - 65);
		}
	}
	return dllRefs;
}

//Run Main Program
func main(){
	file, err := os.Open("dull.in");

	if err != nil {
		log.Fatal(err);
	}

	numDLLs := 0;
	numPrograms := 0;
	numStates := 0;
	bf := bufio.NewReader(file);

	//Loop Until Read a "0" or EOF
	for {
		//Get NPS Values
		line, isPrefix, err := bf.ReadLine();

		if(err != nil || isPrefix) {
			log.Fatal("Problem Reading NPS values");
		}

		if(string(line) == "0"){
			break;
		}

		if err != io.EOF {
			var in []string = strings.Fields(string(line));
			numDLLs, err = strconv.Atoi(in[0]);
			numPrograms, err = strconv.Atoi(in[1]);
			numStates, err = strconv.Atoi(in[2]);
		}

		maxMemory := 0;
		currentMemory := 0;
		progDLL := make(map[int][]int); //Programs With Their DLL's
		locprogMem := make([]int, numPrograms); //Program Memories
		locDLLMem := make([]int, numDLLs); //DLL Memories
		locprogRun := make([]int, numPrograms); //Num Programs Running
		locDLLRun := make([]int, numDLLs);//Num DLLs Running
		locStates := make([]string, numStates); //State Transitions

		//Get DLL Memory Values
		line, isPrefix, err = bf.ReadLine();

		if(err != nil || isPrefix) {
			log.Fatal("Problem loading DLL memory values");
		}

		if err != io.EOF {
			var in []string = strings.Fields(string(line));
			for i := 0; i < len(in); i++ {
				x, err := strconv.Atoi(in[i]);
				if(err != nil){
					log.Fatal("Problem converting DLLMem string to int");
				}
				locDLLMem[i] = x;
			}
		}

		//Load Program Memories and DLL References
		for i := 0; i < numPrograms; i++ {
			line, isPrefix, err = bf.ReadLine();

			if(err != nil || isPrefix) {
				log.Fatal("Problem progMem, progDLL pair");
			}
			var in []string = strings.Fields(string(line));
			mem, err := strconv.Atoi(in[0]);

			if(err != nil){
				log.Fatal(err);
			}
			locprogMem[i] = mem;
			progd := stringConvert(in[1]);

			progDLL[i] = progd;
		}

		//Get The State Transitions (Starting and Stopping Programs)
		line, isPrefix, err = bf.ReadLine();

		if(err != nil || isPrefix) {
			log.Fatal("Problem getting the States");
		}
		locStates = strings.Fields(string(line));

		//Execute States
		//*Really* Ugly Code...
		//for state := 0; state < len(locStates); state++ 
		for state := range locStates{
			progNum, err := strconv.Atoi(locStates[state]);
			if(err != nil){
				log.Fatal("Problem converting state string to int");
			}

			if(progNum > 0 && progNum <= numPrograms){
				//Program Starting
				currentMemory += locprogMem[progNum-1];

				dmem := progDLL[progNum-1];

				for j := 0; j < len(dmem); j++ {
					if(locDLLRun[dmem[j]] == 0){
						currentMemory += locDLLMem[dmem[j]];
					}
					locDLLRun[dmem[j]] += 1;
				}

				if(currentMemory > maxMemory){
					maxMemory = currentMemory;
				}

				locprogRun[progNum-1] += 1;

			} else if(progNum < 0 && (progNum*-1) <= numPrograms){
				//Program Quitting
				if(locprogRun[(progNum*-1)-1] > 0){
					
					currentMemory -= locprogMem[(progNum*-1) - 1];
					dmem := progDLL[(progNum*-1) -1 ];
				
					for j := 0; j < len(dmem); j++ {
						if(locDLLRun[dmem[j]] == 1){
							currentMemory -= locDLLMem[dmem[j]];
						}
						locDLLRun[dmem[j]] -= 1;
					}
					locprogRun[(progNum*-1)-1] -= 1;
				}
			}
		}
		fmt.Println(maxMemory);
	}
}