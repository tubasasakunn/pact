// Pattern 23: State machine with 20+ transitions
component TwentyTransitions {
    states GameCharacter {
        initial Idle

        state Idle { }
        state Walking { }
        state Running { }
        state Jumping { }
        state Falling { }
        state Attacking { }
        state Blocking { }
        state Stunned { }
        state Dead { }

        Idle -> Walking on moveStart
        Idle -> Running on sprint
        Idle -> Jumping on jump
        Idle -> Attacking on attack
        Idle -> Blocking on block
        Walking -> Idle on stop
        Walking -> Running on sprint
        Walking -> Jumping on jump
        Walking -> Attacking on attack
        Running -> Idle on stop
        Running -> Walking on slowDown
        Running -> Jumping on jump
        Running -> Attacking on attack
        Jumping -> Falling on apexReached
        Falling -> Idle on land
        Falling -> Stunned on hardLand
        Attacking -> Idle on attackEnd
        Attacking -> Stunned on interrupted
        Blocking -> Idle on releaseBlock
        Blocking -> Stunned on guardBreak
        Stunned -> Idle on recover
        Stunned -> Dead on fatal
        Dead -> Idle on respawn
    }
}
