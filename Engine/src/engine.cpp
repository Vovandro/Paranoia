//
// Created by devil on 18.05.17.
//

#include "../include/engine.h"

Paranoia::Engine::Engine(eStartType type) {
    run = false;
    this->type = type;

    window = new System::cWindow(this);
    threads = new System::cThreadFactory(this);
    files = new System::cFileFactory(this);
    log = new System::cLog(this, "log");

    log->AddMessage("Init log system", LOG_TYPE::LOG_MESSAGE);

    configs = new Core::cConfigFactory(this);

    render = new Render::cRender(this);
    update = new Core::cUpdate(this);

    states = new Core::cStateManager();

    scenes = new Core::cSceneFactory(this);
}

Paranoia::Engine::~Engine() {
    Stop();

    threads->DestroyFull();

    delete scenes;
    delete states;
    delete update;
    delete render;
    delete configs;
    delete log;
    delete files;
    delete threads;
    delete window;
}


bool Paranoia::Engine::Init(std::string cfName) {

    this->cfName = cfName;

    Core::cConfig *conf = configs->AddObject(cfName, 1, true);

    if (conf == NULL) {
        log->AddMessage("Config file not found: " + cfName, LOG_WARNING);
        log->AddMessage("Load default value from engine", LOG_WARNING);

        conf = new Core::cConfig(this, cfName, 1, true);
    }

    conf->OnAutoCreate();

    window->Init(conf->GetInt("glMajor", 2), conf->GetInt("glMinor", 0), conf->GetInt("antialiasing", 0), false, conf->GetInt("width", 640), conf->GetInt("height", 480), !conf->GetBool("window", true));
    render->Init();



    configs->SaveToFile(cfName);

    run = true;
    return true;
}

void Paranoia::Engine::Start() {
    while (run) {
        handleEvents();

        update->LockLocal();
        render->Update();
        update->UnLockLocal();

        window->Update();

        sf::sleep(sf::milliseconds(1));
    }
}

void Paranoia::Engine::Stop() {
    run = false;
}

void Paranoia::Engine::handleEvents() {
    sf::Event event;

    if (window->GetWindow()->pollEvent(event)) {
        switch (event.type) {
            case sf::Event::Closed:
                run = false;
                break;

            case sf::Event::LostFocus:

                break;

            case sf::Event::GainedFocus:

                break;

            case sf::Event::Resized:
                render->Resize(event.size.width, event.size.height);
                break;

            default:
                break;
        }
    }
}
