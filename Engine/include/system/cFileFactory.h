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
        void Write(std::string fName, cFileData *data);
        void Write(std::string fName, std::string data);
        cFileData *Read(std::string fName, unsigned int size = 0);
        cFileData *ReadLine(std::string fName, unsigned int size = 0);
        cFileData *ReadChar(std::string fName);
        void SetPos(std::string fName, long pos);
        void SetPosStart(std::string fName, long pos);
        void SetPosEnd(std::string fName, long pos);
    };
}

#endif //PROJECT_CFILEFACTORY_H
