#!/bin/bash

curl -vvv -X POST -H "Content-Type: application/json" -d '{"drinks": ["vodka", "gin", "tequila", "rum", "triple sec", "lemon juice"]}' https://cocktails-ewguxkvnaa-uc.a.run.app/cocktails
