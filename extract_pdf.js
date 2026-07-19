ObjC.import('Quartz');

function run(argv) {
    var path = argv[0];
    var url = $.NSURL.fileURLWithPath(path);
    var pdf = $.PDFDocument.alloc.initWithURL(url);
    if (!pdf) {
        console.log("Could not open PDF");
        return;
    }
    
    var numPages = pdf.pageCount;
    for (var i = 0; i < Math.min(numPages, 30); i++) {
        var page = pdf.pageAtIndex(i);
        var text = page.string.js;
        if (text) {
            console.log("--- Page " + (i+1) + " ---");
            console.log(text);
        }
    }
}
