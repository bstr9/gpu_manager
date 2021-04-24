The lifecycle of a container
-----------------------------

## Status
* created
* running
* paused
* restarting
* removing
* exited
* dead

## Status Map
                                                                              [  exited  ]
                                                                                 |   |   
                                                                       docker stop / start             
                                                                                 |   |
             --------                                                            |   v
            | Client |  -> docker create -> [  created  ] -> docker start -> [  started  ] -> [  dead  ] 
             --------                                                         ^  |   ^
                 |                                                            |  |   |
                  -----------------------> docker run ------------------------   |   |
                                                                       docker pause / unpause
                                                                                 |   |
                                                                                 v   |
                                                                             [  paused  ]
