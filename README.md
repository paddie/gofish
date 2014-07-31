gofish
======

gofish is an implementation of a nature article [1] that describes the interaction and movement of a school of fish in 2D.
The package was used to stress-test a checkpointing system, and the interaction model is quite simple:

1. There are two categories of fish: informed and un-informed:
  - Uninformed fish merely react to other fish and their movements.
  - Informed fish have a direction vector that biases them in a certain direction.
2. With every step of the simulation, every fish is attracted to fish that are within a radius RHO (horizon, sight) and..
3. repelled by every fish within a radius ALPHA (personal-sphere, small-radius).

If a certain boundary number of _informed_ fish are introduced into the school (plural of fish, like a murder of crows etc.), the entire school will start moving in the particular direction relative to the average speed of each fish.

[1] Couzin, Iain D., et al. "Effective leadership and decision-making in animal groups on the move." Nature 433.7025 (2005): 513-516.
