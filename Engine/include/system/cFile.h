//
// Created by devil on 18.05.17.
//

#ifndef PROJECT_CFILE_H
#define PROJECT_CFILE_H

#include "../core/cFactoryObject.h"
#include <cstdio>

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
        long size;

        std::string toStr() {return std::string(data);};
    };

    class cFile : public Core::cFactoryObject {
    protected:
        FILE_OPEN_TYPE type;

        std::FILE *file;
        cFileData *fData;

        bool open;
        bool isBinary;

    public:
        cFile(std::string name, int id, bool lock = false);
        virtual ~cFile();

        //Открытие файла
        bool Open(FILE_OPEN_TYPE type, bool binaries = true);
        //Закрытие файла
        void Close();
        //Считывание файла
        cFileData* Read(long size = 0);
        cFileData* ReadLine(long size);
        cFileData* ReadChar();
        //Запись в файл
        void Write(cFileData *data);
        void Write(std::string message);
        //Получение размера файла
        long GetSize();
        //Получение текущего указателя в файле
        long GetPos();
        //Изменение указателя в файле
        void SetPos(long pos);
        //Изменение указателя в файле от начала файла
        void SetPosStart(long pos);
        //Изменение указателя в файле от конца файла
        void SetPosEnd(long pos);
        //Очищает память от считанных данных
        void ClearData();
    };
}

#endif //PROJECT_CFILE_H
