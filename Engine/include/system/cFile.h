//
// Created by devil on 18.05.17.
//

#ifndef PROJECT_CFILE_H
#define PROJECT_CFILE_H

#include "../core/cFactoryObject.h"
#include <fstream>
#include <iostream>
#include <stdio.h>

enum FILE_OPEN_TYPE {
    OPEN_READ,
    OPEN_WRITE,
    OPEN_RW,
};

namespace System {
    class cFile : public Core::cFactoryObject {
    protected:
        FILE_OPEN_TYPE type;
        std::ofstream *f_W;
        std::ifstream *f_R;
        std::fstream *f_RW;

        bool open;

    public:
        cFile(std::string name, int id, bool lock = false);
        ~cFile();

        //Открытие файла
        bool Open(FILE_OPEN_TYPE type);
        //Закрытие файла
        void Close();
        //Считывание файла
        bool Read();
    };
}

#endif //PROJECT_CFILE_H
