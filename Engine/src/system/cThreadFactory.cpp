//
// Created by devil on 18.05.17.
//

#include "../../include/system/cThreadFactory.h"

System::cThreadFactory::cThreadFactory():cFactory() {
}

System::cThreadFactory::~cThreadFactory() {
    for (int i = 0; i <= obj.size(); i++) {
        if (obj[i] != NULL) {
            obj[i]->Destroy();
        }
    }
}

void System::cThreadFactory::AddWork(cThread *work, bool play) {
    AddObject(work);

    work->Init();

    if (play)
        work->Start();
}

void System::cThreadFactory::Play(std::string name) {
    cThread *work = FindObject(name);
    if (work)
        work->Start();
}

void System::cThreadFactory::Pause(std::string name) {
    cThread *work = FindObject(name);
    if (work)
        work->SetEnabled(false);
}

void System::cThreadFactory::Stop(std::string name) {
    cThread *work = FindObject(name);

    if (work != NULL) {
        work->SetEnabled(false);
        work->SetLoop(false);
        work->SetStop();
    }
}

void System::cThreadFactory::Destroy(std::string name) {
    cThread *work = FindObject(name);

    if (work != NULL) {
        work->Destroy();
    }
}

void System::cThreadFactory::Update() {
    for (std::vector<cThread*>::iterator tmp = obj.begin(); tmp != obj.end(); ++tmp) {
        cThread *work = *tmp;
        if (work != NULL) {
            if (work->GetMessage()) {
                work->Message();
            }

            if (work->GetStop()) {
                obj.erase(tmp);
                delete work;
            }
        } else {
            obj.erase(tmp);
        }
    }
}