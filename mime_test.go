package mimetype

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gabriel-vasile/mimetype/internal/matchers"
)

const testDataDir = "testdata"

var files = map[string]*node{
	// archives
	"pdf.pdf":     pdf,
	"zip.zip":     zip,
	"tar.tar":     tar,
	"xls.xls":     xls,
	"xlsx.xlsx":   xlsx,
	"doc.doc":     doc,
	"doc.1.doc":   doc,
	"docx.docx":   docx,
	"docx.1.docx": docx,
	"ppt.ppt":     ppt,
	"pptx.pptx":   pptx,
	"pub.pub":     pub,
	"odt.odt":     odt,
	"ott.ott":     ott,
	"ods.ods":     ods,
	"ots.ots":     ots,
	"odp.odp":     odp,
	"otp.otp":     otp,
	"odg.odg":     odg,
	"otg.otg":     otg,
	"odf.odf":     odf,
	"epub.epub":   epub,
	"7z.7z":       sevenZ,
	"jar.jar":     jar,
	"gz.gz":       gzip,
	"fits.fits":   fits,
	"xar.xar":     xar,
	"bz2.bz2":     bz2,
	"a.a":         ar,
	"deb.deb":     deb,
	"rar.rar":     rar,
	"djvu.djvu":   djvu,
	"mobi.mobi":   mobi,
	"lit.lit":     lit,
	"warc.warc":   warc,
	"zst.zst":     zstd,

	// images
	"png.png":          png,
	"jpg.jpg":          jpg,
	"jp2.jp2":          jp2,
	"jpf.jpf":          jpx,
	"jpm.jpm":          jpm,
	"psd.psd":          psd,
	"webp.webp":        webp,
	"tif.tif":          tiff,
	"ico.ico":          ico,
	"bmp.bmp":          bmp,
	"bpg.bpg":          bpg,
	"heic.single.heic": heic,

	// video
	"mp4.mp4":   mp4,
	"mp4.1.mp4": mp4,
	"webm.webm": webM,
	"3gp.3gp":   threeGP,
	"3g2.3g2":   threeG2,
	"flv.flv":   flv,
	"avi.avi":   avi,
	"mov.mov":   quickTime,
	"mqv.mqv":   mqv,
	"mpeg.mpeg": mpeg,
	"mkv.mkv":   mkv,
	"asf.asf":   asf,

	// audio
	"mp3.mp3":            mp3,
	"mp3.v1.notag.mp3":   mp3,
	"mp3.v2.notag.mp3":   mp3,
	"mp3.v2.5.notag.mp3": mp3,
	"wav.wav":            wav,
	"flac.flac":          flac,
	"midi.midi":          midi,
	"ape.ape":            ape,
	"aiff.aiff":          aiff,
	"au.au":              au,
	"ogg.oga":            oggAudio,
	"ogg.spx.oga":        oggAudio,
	"ogg.ogv":            oggVideo,
	"amr.amr":            amr,
	"mpc.mpc":            musePack,
	"aac.aac":            aac,
	"voc.voc":            voc,
	"m4a.m4a":            m4a,
	"m4b.m4b":            aMp4,
	"qcp.qcp":            qcp,

	// source code
	"html.html":         html,
	"html.withbr.html":  html,
	"svg.svg":           svg,
	"svg.1.svg":         svg,
	"txt.txt":           txt,
	"php.php":           php,
	"ps.ps":             ps,
	"json.json":         json,
	"geojson.geojson":   geoJson,
	"geojson.1.geojson": geoJson,
	"ndjson.ndjson":     ndJson,
	"csv.csv":           csv,
	"tsv.tsv":           tsv,
	"rtf.rtf":           rtf,
	"js.js":             js,
	"lua.lua":           lua,
	"pl.pl":             perl,
	"py.py":             python,
	"tcl.tcl":           tcl,
	"vCard.vCard":       vCard,
	"vCard.dos.vCard":   vCard,
	"ics.ics":           iCalendar,
	"ics.dos.ics":       iCalendar,

	// binary
	"class.class": class,
	"swf.swf":     swf,
	"crx.crx":     crx,
	"wasm.wasm":   wasm,
	"exe.exe":     exe,
	"ln":          elfExe,
	"so.so":       elfLib,
	"o.o":         elfObj,
	"dcm.dcm":     dcm,
	"mach.o":      macho,
	"sample32":    macho,
	"sample64":    macho,
	"mrc.mrc":     mrc,

	// fonts
	"woff.woff":   woff,
	"woff2.woff2": woff2,
	"otf.otf":     otf,
	"eot.eot":     eot,

	// XML and subtypes of XML
	"xml.withbr.xml": xml,
	"kml.kml":        kml,
	"xlf.xlf":        xliff,
	"dae.dae":        collada,
	"gml.gml":        gml,
	"gpx.gpx":        gpx,
	"tcx.tcx":        tcx,
	"x3d.x3d":        x3d,
	"amf.amf":        amf,
	"3mf.3mf":        threemf,
	"rss.rss":        rss,
	"atom.atom":      atom,

	"shp.shp": shp,
	"shx.shx": shx,
	"dbf.dbf": dbf,

	"sqlite3.sqlite3": sqlite3,
	"dwg.dwg":         dwg,
	"dwg.1.dwg":       dwg,
	"nes.nes":         nes,
	"mdb.mdb":         mdb,
	"accdb.accdb":     accdb,
}

func TestMatching(t *testing.T) {
	errStr := "File: %s; Mime: %s != DetectedMime: %s; err: %v"
	for fName, node := range files {
		fileName := filepath.Join(testDataDir, fName)
		f, err := os.Open(fileName)
		if err != nil {
			t.Fatal(err)
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}

		if dMime, _ := Detect(data); dMime != node.mime {
			t.Errorf(errStr, fName, node.mime, dMime, nil)
		}

		if _, err := f.Seek(0, io.SeekStart); err != nil {
			t.Errorf(errStr, fName, node.mime, root.mime, err)
		}

		if dMime, _, err := DetectReader(f); dMime != node.mime {
			t.Errorf(errStr, fName, node.mime, dMime, err)
		}
		f.Close()

		if dMime, _, err := DetectFile(fileName); dMime != node.mime {
			t.Errorf(errStr, fName, node.mime, dMime, err)
		}
	}
}

func TestFaultyInput(t *testing.T) {
	inexistent := "inexistent.file"
	if _, _, err := DetectFile(inexistent); err == nil {
		t.Errorf("%s should not match successfully", inexistent)
	}

	f, _ := os.Open(inexistent)
	if _, _, err := DetectReader(f); err == nil {
		t.Errorf("%s reader should not match successfully", inexistent)
	}
}

func TestEmptyInput(t *testing.T) {
	if m, _ := Detect([]byte{}); m != "inode/x-empty" {
		t.Errorf("failed to detect empty file")
	}
}

func TestBadBdfInput(t *testing.T) {
	if m, _, _ := DetectFile("testdata/bad.dbf"); m != "application/octet-stream" {
		t.Errorf("failed to detect bad DBF file")
	}
}

func TestGenerateSupportedMimesFile(t *testing.T) {
	f, err := os.OpenFile("supported_mimes.md", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	nodes := root.flatten()
	header := fmt.Sprintf(`## %d Supported MIME types
This file is automatically generated when running tests. Do not edit manually.

Extension | MIME type
--------- | --------
`, len(nodes))

	if _, err := f.WriteString(header); err != nil {
		t.Fatal(err)
	}
	for _, n := range nodes {
		ext := n.extension
		if ext == "" {
			ext = "n/a"
		}
		str := fmt.Sprintf("**%s** | %s\n", ext, n.mime)
		if _, err := f.WriteString(str); err != nil {
			t.Fatal(err)
		}
	}
}

func TestIndexOutOfRange(t *testing.T) {
	for _, n := range root.flatten() {
		_ = n.matchFunc(nil)
	}
}

func BenchmarkMatchDetect(b *testing.B) {
	files := []string{"png.png", "jpg.jpg", "pdf.pdf", "zip.zip", "docx.docx", "doc.doc"}
	data, fLen := [][matchers.ReadLimit]byte{}, len(files)
	for _, f := range files {
		d := [matchers.ReadLimit]byte{}

		file, err := os.Open(filepath.Join(testDataDir, f))
		if err != nil {
			b.Fatal(err)
		}

		io.ReadFull(file, d[:])
		data = append(data, d)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Detect(data[n%fLen][:])
	}
}
