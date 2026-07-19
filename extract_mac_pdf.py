import sys
import CoreGraphics as cg

def extract_pdf_text(pdf_path):
    pdf_url = cg.CGDataProviderCreateWithFilename(pdf_path)
    if pdf_url is None:
        return "Could not open file"
    pdf = cg.CGPDFDocumentCreateWithProvider(pdf_url)
    if pdf is None:
        return "Could not open PDF"
    
    num_pages = pdf.getNumberOfPages()
    for i in range(1, min(num_pages + 1, 30)):
        page = pdf.getPage(i)
        # CoreGraphics in python 2.7 had methods, but we are in python3.
        # Python 3 doesn't have CoreGraphics.
        pass

extract_pdf_text('bahan/SN_SRS_v5.pdf')
