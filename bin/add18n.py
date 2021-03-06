import json
import os
import sys

import pyperclip

controller = sys.argv[1]
key = sys.argv[2].replace(" ","").replace("/","").replace("\\","").replace("`","").replace("~","").replace("@","").replace("#","").replace("$","").replace("%","").replace("^","").replace("&","").replace("*","").replace("(","").replace(")","").replace("-","").replace("+","").replace("=","").replace("[","").replace("]","").replace("{","").replace("}","").replace("|","").replace(":","").replace(";","").replace("'","").replace("\"","").replace("<","").replace(">","").replace(",","").replace(".","").replace("?","").replace("/","").replace("!","")
key = key[:35]

translation = sys.argv[3]

blastIt = sys.argv[4]

if translation[:1] == "{":
    translation = translation[1:]

if translation[len(translation)-1:len(translation)] == "}":
    translation = translation[:len(translation)-1]
gopath = os.getenv("GOPATH")
controllers = "".join(open(gopath + '/src/github.com/blauentag/blau/controllers/constants.go', 'r').readlines())

app = "app"
page = "app"
if controllers.find(controller) > 0 and controller != "app":
    app = controller
    page = "page"
    controller = ""


path = '/src/github.com/blauentag/blau/web/app/globalization/translations/' + app + '/en/US.json'
translations = json.load(open(gopath + path, 'r'))
fullString = translation.strip()
if blastIt == "0" and key not in translations:
    translations[key] = fullString
    json.dump(translations, open(gopath + path, 'w'), indent=4, sort_keys=True)
    jsCode = 'window.' + page + 'Content.' + key
    print 'floatingLabelText={' + jsCode + '}'
    try:
        golang=False
        if golang:
            pyperclip.copy("queries.AppContent.GetTranslation(context, \"" + key + "\")")
        else:
            pyperclip.copy("{" + jsCode + "}")
    except:
        pass

    if app == "app":
        translationGoCode = gopath + '/src/github.com/blauentag/blau/queries/appTranslations.go'
        with open(translationGoCode, 'r') as content_file:
            appContent = content_file.read()
        if key + " " not in appContent:
            f = open(translationGoCode,'w')
            appContent = appContent.replace("//AdditionalConstructs", "\t" + key + "                            string `json:\"" + key + "\"`\n\t//AdditionalConstructs")
            f.write(appContent)
            f.close()
else:
    try:
	    jsCode = 'window.' + page + 'Content.' + key
	    pyperclip.copy("{" + jsCode + "}")
    except:
        pass
    print "TranslationExists"
