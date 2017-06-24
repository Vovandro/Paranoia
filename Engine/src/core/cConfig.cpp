//
// Created by devil on 01.06.17.
//

#include "../../include/core/cConfig.h"
#include "../../include/engine.h"

Core::cConfig::cConfig(Paranoia::Engine *engine, std::string name, int id, bool lock) : Core::cFactoryObject(engine, name, id, lock) {
    AutoValue = false;
}

std::string Core::cConfig::GetLine(std::string *text, std::string split, int &start) {
    int pos = text->find(split, start);
    std::string ret = text->substr(start, pos - start);
    start = pos + split.size();
    return ret;
}

std::string Core::cConfig::ToString() {
    std::string ret;
    ret += std::to_string(items.size());

    for (int i = 0; i < items.size(); i++) {
        ret +=  "|?|" + std::to_string(i) + "|||" + items[i]->name + "=";

        if (items[i]->type == CIT_STRING) {
            ret += std::to_string(CIT_STRING) + "==" + ((cConfigItemString*)items[i])->data;
        }

        if (items[i]->type == CIT_FLOAT) {
            ret += std::to_string(CIT_FLOAT) + "==" + std::to_string(((cConfigItemFloat*)items[i])->data);
        }

        if (items[i]->type == CIT_INT) {
            ret += std::to_string(CIT_INT) + "==" + std::to_string(((cConfigItemInt*)items[i])->data);
        }
    }

    return ret;
}

void Core::cConfig::FromString(std::string text) {
    int j = 0;
    int size = std::stoi(GetLine(&text, "|?|", j));

    for (int i = 0; i < size; i++) {
        if (j < text.size()) {
            int num = std::stoi(GetLine(&text, "|||", j));

            if (num != i)
                return;

            std::string name = GetLine(&text, "=", j);
            CONFIG_ITEM_TYPE type = (CONFIG_ITEM_TYPE) std::stoi(GetLine(&text, "==", j));

            switch (type) {
                case CIT_STRING: {
                    cConfigItemString *newItem = new cConfigItemString();

                    newItem->name = name;
                    newItem->data = GetLine(&text, "|?|", j);

                    items.push_back((cConfigItem*)newItem);
                }
                break;

                case CIT_FLOAT: {
                    cConfigItemFloat *newItem = new cConfigItemFloat();

                    newItem->name = name;
                    newItem->data = std::stof(GetLine(&text, "|?|", j));

                    items.push_back((cConfigItem*)newItem);
                }
                break;

                case CIT_INT: {
                    cConfigItemInt *newItem = new cConfigItemInt();

                    newItem->name = name;
                    newItem->data = std::stoi(GetLine(&text, "|?|", j));

                    items.push_back((cConfigItem*)newItem);
                }
                break;

                default:
                    break;
            }
        }
    }
}

void Core::cConfig::Add(Core::cConfigItem *item) {
    items.push_back(item);
}

void Core::cConfig::Register() {
    if (engine)
    {
        engine->configs->AddObject(this);
    }
}

void Core::cConfig::OnAutoCreate() {
    AutoValue = true;
}

std::string Core::cConfig::GetString(std::string vName, std::string def) {
    for (int i = 0; i < items.size(); i++) {
        if ((items[i]->name == vName)&&(items[i]->type == CIT_STRING)) {
            return ((cConfigItemString*)items[i])->data;
        }
    }

    if (AutoValue) {
        cConfigItemString *newItem = new cConfigItemString();
        newItem->name = vName;
        newItem->data = def;

        items.push_back((cConfigItem*)newItem);
    }

    return def;
}

float Core::cConfig::GetFloat(std::string vName, float def) {
    for (int i = 0; i < items.size(); i++) {
        if ((items[i]->name == vName)&&(items[i]->type == CIT_FLOAT)) {
            return ((cConfigItemFloat*)items[i])->data;
        }
    }

    if (AutoValue) {
        cConfigItemFloat *newItem = new cConfigItemFloat();
        newItem->name = vName;
        newItem->data = def;

        items.push_back((cConfigItem*)newItem);
    }

    return def;
}

int Core::cConfig::GetInt(std::string vName, int def) {
    for (int i = 0; i < items.size(); i++) {
        if ((items[i]->name == vName)&&(items[i]->type == CIT_INT)) {
            return ((cConfigItemInt*)items[i])->data;
        }
    }

    if (AutoValue) {
        cConfigItemInt *newItem = new cConfigItemInt();
        newItem->name = vName;
        newItem->data = def;

        items.push_back((cConfigItem*)newItem);
    }

    return def;
}
