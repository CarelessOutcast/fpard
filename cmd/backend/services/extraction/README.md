
# Extracting content from PDF

I originally wanted to create a pure go implementation of pdf extraction but I
wanted to finish the initial prototype quickly. I figured that I could use
some sort of pdf extraction software in go but I didn't find one that
completely fit my needs. For the inital version of this, I'm going to use 
python and incorporate the library into this. 

I'm going to be making a binary and embedding it into the final go binary. 

## Building the binary
``` powershell
pyinstaller --onefile 

```

