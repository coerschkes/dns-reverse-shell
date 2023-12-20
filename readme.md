<h3 align="center">dns reverse shell</h3>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->

## About The Project

This project is based on the idea of having a simple reverse-tcp shell for educational purposes but instead of
just reversing tcp it uses dns requests to communicate with the server to hide the actual communication from IDS/IPS
systems.
This is done by submitting the payload and commands as dns queries. To prevent a dns header overflow the payload is
split into multiple queries.

## Getting Started

todo: add instructions
<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage

todo: add usage
<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Roadmap

- [X] Simple navigation
- [X] Dns communication
- [X] Message splitter
- [X] Polling (if idle only poll twice per minute to prevent flooding)
- [ ] Protocol for sending "big" messages -> currently server receives no answer because the client does not send anything
- [ ] Navigation improvement: navigation relative to the current path
- [ ] Command chaining
- [ ] Multiple sessions
- [ ] Windows navigation support
- [ ] Use encryption for payload
- [ ] Autostart

<p align="right">(<a href="#readme-top">back to top</a>)</p>