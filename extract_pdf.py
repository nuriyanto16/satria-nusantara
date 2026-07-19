import sys, zlib, re

def extract(pdf_path):
    with open(pdf_path, 'rb') as f:
        pdf_data = f.read()

    # Find all streams
    streams = re.findall(b'stream[\r\n]+(.*?)[\r\n]+endstream', pdf_data, re.DOTALL)
    
    text_content = []
    for s in streams:
        try:
            # Try to decompress
            decompressed = zlib.decompress(s)
            # Find text within Tj or TJ operators roughly
            # This is very basic and won't handle complex encodings or CMap
            texts = re.findall(b'\((.*?)\) Tj', decompressed)
            texts += re.findall(b'\((.*?)\) TJ', decompressed)
            for t in texts:
                try:
                    text_content.append(t.decode('utf-8', 'ignore'))
                except:
                    pass
        except Exception as e:
            pass

    print('\n'.join(text_content))

extract('bahan/SN_SRS_v5.pdf')
