import pymupdf
import sys

def extractText(filename):
    doc = pymupdf.open(filename)

    output = bytearray()
    for page in doc:
        output.extend(page.get_text().encode("utf8"))
        print("[" + ",".join(str(b) for b in output) + "]")


def main():
    if len(sys.argv) != 2:
        print("Usage: python script.py <filename>", file=sys.stderr)
        sys.exit(1)

    filename = sys.argv[1]
    try: 
        extractText(filename)
    except:
        sys.exit(1)


if __name__ == "__main__":
    main()
