package main
//test
import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"encoding/csv"
	"github.com/sqweek/dialog"
)


func main(){
	filename := loadfile()
	BPRfolder :=prepareBPR(filename)

	folder := filepath.Dir(filename)
	logFilePath:= filepath.Join(folder,"log.txt")
os.Create(logFilePath)
logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {log.Println(err.Error())}
fmt.Fprintln(logFile,"Logfile lmao")
log.SetOutput(logFile)
	tempFolder := filepath.Join(folder,"temp")
	manifestPath := filepath.Join(tempFolder,"manifest.xml")

	attachmentFolder := extractFiles(filename)
	Bprfiles := extractInfo(manifestPath)
	matchKey := generateKey("csv")
	unmatchedFolder := filepath.Join(folder,"unmatched")
	err = os.Mkdir(unmatchedFolder,os.ModePerm)
	if err != nil {log.Println(err.Error())}
fmt.Println("created unmatched folder")


	Bprfiles = matchFiles(Bprfiles,matchKey) //this updates the Bprfiles
	//Create an index of files to be copied
	var copyIndex []transferInstructions

	for i :=0; i<len(Bprfiles);i++{

		if Bprfiles[i].Matched == true{
		var index transferInstructions
		index.from = filepath.Join(attachmentFolder,Bprfiles[i].MD5) 
		index.to = filepath.Join(BPRfolder,Bprfiles[i].BPRFolder,Bprfiles[i].RealName)
		copyIndex = append(copyIndex, index)
	}}
	copyFromIndex(copyIndex,false)	
	//Create an index of files that were not matched.

	var unmatchedIndex []transferInstructions
	unmatched := 0
	for i :=0; i<len(Bprfiles);i++{
		if Bprfiles[i].Matched == false{
			var index transferInstructions
			index.from = filepath.Join(attachmentFolder,Bprfiles[i].MD5) 
			index.to = filepath.Join(unmatchedFolder,Bprfiles[i].subtype,Bprfiles[i].RealName)
			unmatchedSubFolder := filepath.Dir(index.to)
	err := os.Mkdir(unmatchedSubFolder,os.ModePerm)
	if err != nil {log.Println(err.Error())}
			unmatchedIndex = append(unmatchedIndex, index)
			unmatched++ 
		}
	}
	if unmatched >0{

 fmt.Println("Moving ",unmatched," unmatched files.")
		copyFromIndex(unmatchedIndex,true)
	}

	os.RemoveAll(tempFolder)

defer logFile.Close()

}





func loadfile()(filename string ){
	filename, err := dialog.File().Filter("IUCLID 6 dossier", "i6z").Load()	
	if err != nil {log.Println(err.Error())}
	filename = filepath.Clean(filename)
	folder := filepath.Dir(filename)
	fmt.Println("Working directory:",folder)
	return
}
func prepareBPR(i6zPath string)(bprFolderPath string ){
	bprFolders := []string{"1 Applicant","2 Identity of the Biocidal Product","3 Physical, chemical and technical properties", "4 Physical hazards and respective characteristics", "5 Methods of detection and identification","6 Effectiveness against target organisms","7 Intended uses and exposure", "8 Toxicological profile for humans and animals","9 Ecotoxicological studies","10 Environmental fate and behaviour","11 Measures to protect humans, animals and the environment","12 Classification and labelling","13 Summary and evaluation"}

	Folder := filepath.Dir(i6zPath)
	bprFolderPath = filepath.Join(Folder,"BPR")
	err := os.Mkdir(bprFolderPath,os.ModePerm)
	if err != nil {log.Println(err.Error())}

	for _, f :=range bprFolders{
		subfolder := filepath.Join(bprFolderPath,f)
		err = os.Mkdir(subfolder,os.ModePerm)
	}
	fmt.Println("BPR Folders prepared")
	return
}
func extractFiles(i6zPath string)(attachmentFolder string){

	folder := filepath.Dir(i6zPath)
	tempFolderPath := filepath.Join(folder,"temp")
	attachmentFolder = filepath.Join(tempFolderPath,"attachments")
	i6z,err := zip.OpenReader(i6zPath)
	if err != nil {log.Print(err.Error())}
	defer	i6z.Close()
	err = os.Mkdir(tempFolderPath,os.ModePerm)
	err = os.Mkdir(tempFolderPath+"/attachments",os.ModePerm)

	for _, f := range i6z.File {
		extractPath := filepath.Join(tempFolderPath,f.Name)
		// fmt.Printf("Extracting file:: %s\n", extractPath)
		//
		outFile, err := os.OpenFile(extractPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
	if err != nil {log.Println(err.Error())}
		}
		fileInArchive, err := f.Open()
		if err != nil {
	if err != nil {log.Println(err.Error())}
		}
		if _, err := io.Copy(outFile, fileInArchive); err != nil {
	if err != nil {log.Println(err.Error())}
		}

		outFile.Close()
		fileInArchive.Close()	
	}
	fmt.Println("Files extracted")
	return
}


func extractInfo(manifestPath string)(files []fileInfo){

	manifest,err := os.Open(manifestPath)
	if err != nil {log.Println(err.Error())}
	fmt.Println("Opened manifest.xml")
	defer manifest.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(manifest)
	// we initialize our Users array
	var dossier dossier
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	err = xml.Unmarshal(byteValue, &dossier)
	if err != nil {log.Println(err.Error())}
	var Files []fileInfo

	for j, f := range dossier.Attachment {
		for i, g := range dossier.Document{
			if  f.Container == g.Container {
				var file fileInfo
				file.UUID = f.Container
				file.subtype = g.Category
				file.MD5 = filepath.Base(f.MD5Filename.LinkedDoc)
				file.RealName = f.RealFilename
				Files = append(Files, file)
				i++
			} 
		}

		j++
	}
	fmt.Println("Dossier info extracted from XML")
	return Files
}


func generateKey(csvFile string)(CSVkey []legislationKey){
	if csvFile == "default"{
		fmt.Println("Using default BPR definitions")
	bprFolders := []string{"1 Applicant","2 Identity of the Biocidal Product","3 Physical, chemical and technical properties", "4 Physical hazards and respective characteristics", "5 Methods of detection and identification","6 Effectiveness against target organisms","7 Intended uses and exposure", "8 Toxicological profile for humans and animals","9 Ecotoxicological studies","10 Environmental fate and behaviour","11 Measures to protect humans, animals and the environment","12 Classification and labelling","13 Summary and evaluation"}

	bprTypes := []string{"Sites","Identifiers","GeneralInformation","Flammability","AnalyticalMethods","EffectivenessAgainstTargetOrganisms","BioIntendedUsesExposure","ToxicologicalProfile","EcotoxicologicalInformation","EnvironmentalFateAndPathways","ProtectionMeasures","Ghs","BioSummaryEvaluation"}

	var defaultKey []legislationKey
	for i:=0; i<len(bprFolders); i++{
		var key legislationKey
		key.XMLkey=bprTypes[i]
		key.section=bprFolders[i]
	defaultKey = append(defaultKey, key)
	}
		return defaultKey
	} else {

	IUCdefinition, err := dialog.File().Filter( ".csv").Load()	
	if err != nil {log.Println(err.Error())}
	file,err := os.OpenFile(IUCdefinition,os.O_RDONLY,0444)
	if err != nil {log.Println(err.Error())}
	defer file.Close()
	raw:= csv.NewReader(file)
	raw.LazyQuotes = true
	raw.FieldsPerRecord = -1
	output,err := raw.ReadAll()
	if err != nil {
	if err != nil {log.Println(err.Error())}
	}   

	var CSVkey []legislationKey

	for r := range output{
		var oneKey legislationKey
		fields :=(len(output[r]))
		if fields>2{
			for i:=1; i<=fields-1;i++{
				oneKey.section = output[r][0]
				oneKey.XMLkey= output[r][i]

				CSVkey = append(CSVkey, oneKey)
			}
		}else{
			oneKey.section = output[r][0]
			oneKey.XMLkey= output[r][1]
			CSVkey = append(CSVkey, oneKey)
		}
	}
		return CSVkey
	}
}
func matchFiles(Bprfiles []fileInfo, CSVkey []legislationKey)(MatchedFiles []fileInfo){


	numberOfFiles := len(Bprfiles)
	numberOfMatches := 0 

	for j := range Bprfiles{
		Bprfiles[j].Matched = false
		for k := range CSVkey{
			if Bprfiles[j].subtype == CSVkey[k].XMLkey {
				Bprfiles[j].Matched = true
				numberOfMatches++
				Bprfiles[j].BPRFolder = CSVkey[k].section
			} 
		}
	}
	fmt.Println("A total of ",numberOfFiles," files produced ",numberOfMatches," matches. Files unaccounted for:",numberOfFiles - numberOfMatches)
	MatchedFiles = Bprfiles
	return MatchedFiles
}

func copyFromIndex(copyIndex []transferInstructions, remove bool)(){
for _, transfer := range copyIndex{


		src := transfer.from
		dst := transfer.to 

		fin, err := os.Open(src)
	if err != nil {log.Println(err.Error())}
		defer fin.Close()

		fout, err := os.Create(dst)
		defer fout.Close()

	if err != nil {log.Println(err.Error())}
		_, err = io.Copy(fout, fin)

	if err != nil {log.Println(err.Error())}
	}
	if remove == true{
	for d := range copyIndex{
	os.Remove(copyIndex[d].from)	}
}}

func matchCustomCSV(Bprfiles []fileInfo)(MatchedFiles []fileInfo){
	IUCdefinition, err := dialog.File().Filter( ".csv").Load()	
	if err != nil {log.Println(err.Error())}
	file,err := os.OpenFile(IUCdefinition,os.O_RDONLY,0444)
	if err != nil {log.Println(err.Error())}
	defer file.Close()
	raw:= csv.NewReader(file)
	raw.LazyQuotes = true
	raw.FieldsPerRecord = -1
	output,err := raw.ReadAll()

	if err != nil {log.Println(err.Error())}
	var CSVkey []legislationKey

	for r := range output{
		var oneKey legislationKey
		fields :=(len(output[r]))
		if fields>2{fmt.Println(fields)
			for i:=1; i<=fields-1;i++{
				oneKey.section = output[r][0]
				oneKey.XMLkey= output[r][i]

				CSVkey = append(CSVkey, oneKey)
			}
		}else{
			oneKey.section = output[r][0]
			oneKey.XMLkey= output[r][1]
			CSVkey = append(CSVkey, oneKey)
		}
	}
	numberOfFiles := len(Bprfiles)
	numberOfMatches := 0 

	for j := range Bprfiles{
		Bprfiles[j].Matched = false
		for k := range CSVkey{
			if Bprfiles[j].subtype == CSVkey[k].XMLkey {
				Bprfiles[j].Matched = true
				numberOfMatches++
				Bprfiles[j].BPRFolder = CSVkey[k].section
			} 
		}
	}
	fmt.Println("A total of ",numberOfFiles," files produced ",numberOfMatches," matches. Files unaccounted for:",numberOfFiles - numberOfMatches)
	MatchedFiles = Bprfiles
	return MatchedFiles
}




