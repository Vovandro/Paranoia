//
// Created by devil on 01.06.17.
//

#ifndef PROJECT_CCONFIG_H
#define PROJECT_CCONFIG_H

#include "vector"

namespace Core {
    class cConfigItem {
    public:
        std::string name;
    };

    class cConfigItemString : public cConfigItem {
    public:
        std::string data;
    };

    class cConfigItemFloat : public cConfigItem {
    public:
        float data;
    };

    /* Класс для работы с конфигурациями для всех объектов записанных в файлах */
    class cConfig {
    protected:
        std::vector<cConfigItem*> items;

    public:
        std::string ToString();
        void FromString(std::string text);
    };
}

#endif //PROJECT_CCONFIG_H
