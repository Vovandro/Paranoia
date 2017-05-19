//
// Created by devil on 18.05.17.
//

#ifndef PROJECT_CFILEFACTORY_H
#define PROJECT_CFILEFACTORY_H

#include "cFile.h"
#include "../core/cFactory.h"

namespace System {
    class cFileFactory : public Core::cFactory<cFile> {
    protected:
    public:
        bool Add(std::string fName, FILE_OPEN_TYPE type);
        cFileData *Read(std::string fName, unsigned int size, bool isLine = false, bool isWord = false);
    };
}

#endif //PROJECT_CFILEFACTORY_H
