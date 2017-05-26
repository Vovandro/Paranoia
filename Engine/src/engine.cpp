//
// Created by devil on 18.05.17.
//

#include "../include/engine.h"

Paranoia::Engine::Engine(eStartType type) {
    run = false;
    this->type = type;

    window = new System::cWindow(this);
    threads = new System::cThreadFactory();
    files = new System::cFileFactory();
    log = new System::cLog(this, "log");

    log->AddMessage("Init log system", LOG_TYPE::LOG_MESSAGE);

    render = new Render::cRender(this);
    update = new Core::cUpdate(this);

    states = new Core::cStateManager();
}

Paranoia::Engine::~Engine() {
    Stop();

    delete states;
    delete update;
    delete render;
    delete log;
    delete files;
    delete threads;
    delete window;
}


bool Paranoia::Engine::Init() {

    window->Init(2, 2, 0);
    render->Init();

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
