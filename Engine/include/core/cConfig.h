//
// Created by devil on 01.06.17.
//

#ifndef PROJECT_CCONFIG_H
#define PROJECT_CCONFIG_H

#include "vector"
#include "cFactoryObject.h"

enum CONFIG_ITEM_TYPE {
    CIT_STRING,
    CIT_INT,
    CIT_FLOAT,
};

namespace Core {
    class cConfigItem {
    public:
        std::string name;
        CONFIG_ITEM_TYPE type;
    };

    class cConfigItemString : public cConfigItem {
    public:
        cConfigItemString() {type = CIT_STRING;};
        std::string data;
    };

    class cConfigItemFloat : public cConfigItem {
    public:
        cConfigItemFloat() {type = CIT_FLOAT;};
        float data;
    };

    class cConfigItemInt : public cConfigItem {
    public:
        cConfigItemInt() {type = CIT_INT;};
        int data;
    };

    /* Класс для работы с конфигурациями для всех объектов записанных в файлах */
    class cConfig : public cFactoryObject {
    protected:
        std::vector<cConfigItem*> items;
        std::string GetLine(std::string *text, std::string split, int &start);

    public:
        cConfig(std::string name, int id, bool lock = false);

        std::string ToString();
        void FromString(std::string text);

        void Add(cConfigItem* item);
    };
}

#endif //PROJECT_CCONFIG_H
