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
    CIT_BOOL,
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

    class cConfigItemBool : public cConfigItem {
    public:
        cConfigItemBool() {type = CIT_BOOL;};
        bool data;
    };

    /* Класс для работы с конфигурациями для всех объектов записанных в файлах */
    class cConfig : public cFactoryObject {
    protected:
        std::vector<cConfigItem*> items;
        std::string GetLine(std::string *text, std::string split, int &start);
        bool AutoValue;

    public:
        cConfig(Paranoia::Engine *engine, std::string name, int id, bool lock = false);

        virtual void Register() override;

        // Собирает все свои параметры в строку
        std::string ToString();
        // Разбивает строку на параметры
        void FromString(std::string text);

        // Добавление нового параметра
        void Add(cConfigItem* item);

        // Включение режима автоматического создания переменной при ее отсутствии с дефолтным параметром
        void OnAutoCreate();

        std::string GetString(std::string vName, std::string def);
        float GetFloat(std::string vName, float def);
        int GetInt(std::string vName, int def);
        bool GetBool(std::string vName, bool def);
    };
}

#endif //PROJECT_CCONFIG_H
