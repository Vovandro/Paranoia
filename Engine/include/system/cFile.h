//
// Created by devil on 18.05.17.
//

#ifndef PROJECT_CFILE_H
#define PROJECT_CFILE_H

#include "../core/cFactoryObject.h"
#include <stdio.h>

enum FILE_OPEN_TYPE {
    OPEN_READ,
    OPEN_WRITE,
    OPEN_WRITE_CLEAR,
    OPEN_RW,
    OPEN_RW_CLEAR,
};

namespace System {
    class cFileData {
    public:
        char* data;
        int size;
    };

    class cFile : public Core::cFactoryObject {
    protected:
        FILE_OPEN_TYPE type;

        std::FILE *file;

        bool open;

    public:
        cFile(std::string name, int id, bool lock = false);
        virtual ~cFile();

        //Открытие файла
        bool Open(FILE_OPEN_TYPE type, bool binaries = false);
        //Закрытие файла
        void Close();
        //Считывание файла
        cFileData* Read(unsigned int size, bool isLine = false, bool isWord = false);
        //Запись в файл
        void Write(cFileData *data);
        void Write(std::string message);
    };
}

#endif //PROJECT_CFILE_H
