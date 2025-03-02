# Go-FParD

The functional PDF audio reading doohickey

The goal of this project is to create an application that reads PDFs out loud.
I want something shiny that doesn't come with a shiny price. The goal is to
self-host this service, and make it simple enough for anyone to use it. 

---

As it currently stands, I would need to recompile the program and pyinstaller
everytime that I need to update. Maybe creating a rest api for the python
aspect might be more beneficial. I would be able to update and release it's
own modules independently. 

---

I'm also incorporating a C library into the program: Espeak-ng. Since this is
a backend service, I want to output raw PCM audio (with a Wav header) to the
frontend. I don't want to create an intermediate file, and stream the
extracted synthesized text directly to the frontend. I will need to
incorporate the DLL into the build so that this would function. 

