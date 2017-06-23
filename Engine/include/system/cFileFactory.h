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
        cFileFactory(Paranoia::Engine *engine);
        ~cFileFactory();
        // Открытие файла
        virtual bool Open(std::string fName, FILE_OPEN_TYPE type);
        // Запись данных в файл
        void Write(std::string fName, cFileData *data);
        void Write(std::string fName, std::string data);
        // Считывание из файла данных
        cFileData *Read(std::string fName, unsigned int size = 0);
        cFileData *ReadLine(std::string fName, unsigned int size = 0);
        cFileData *ReadChar(std::string fName);
        // Изменение положени курсора в файле
        void SetPos(std::string fName, long pos);
        void SetPosStart(std::string fName, long pos);
        void SetPosEnd(std::string fName, long pos);
    };
}

#endif //PROJECT_CFILEFACTORY_H
