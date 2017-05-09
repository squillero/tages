Payoff
======

Let **R** (*reward*) be the payoff for a mutual cooperation, and **P** (*punishment*) the payoff if both players defect; when only one player acts selfishly, its payoff is **T** (*temptation*), and the payoff of its opponent is **S** (*sucker*). The Prisoner's Dilemma model requires that (**T** > **R** > **P** > **S**) and (2 x **R** > **S** + **T**).

The most common payoff matrix found in literature sets **T**=5, **R**=3, **P**=1, and **S**=0.

Standard Strategies
===================

**2TFT** (*Two Tits For Tat*) -- Cooperates unless the opponent defects. To retort, it defects twice. Then, if the opponent cooperated, it starts cooperating again.

**AC** (*Always Cooperate*) -- Always cooperates.

**AD** (*Always Defect*) -- Always defects.

**ADP** (*Adaptive*) -- Starts with a sequence of five consecutive cooperations, followed by five defections; from the 11-th turn, it chooses the move that has been most profitable so far.

**APAV** (*Adaptive Pavlov*) -- Categorizes the opponent's behavior using four classes, and then plays accordingly. If the opponent is *fully cooperative*, APAV behaves like a simple *TFT*; if it is *almost cooperative*, APAV plays *Soft Tit for Two Tats* in order to recover mutual cooperation; if *aggressive* or *random*, APAV always defects. To categorize the opponent, APAV plays six turns as *TFT*. If the opponent started with a cooperation, it is identified as *fully cooperative*; if the opponent defected three times, it is identified as *almost cooperative*; if the opponent defected four or more times, it is identified as *aggressive*; in all other cases, the opponent is considered *random*. In order to deal with the situations in which the opponents may change their actions, the average payoff is computed every six turns. If it is lower than R, the process of opponent identification is restarted.

**ATFT** (*Adaptive Tit for Tat*) -- Updates the variable w in [0,1] according to the opponent's moves: w is slowly pushed toward 1 on cooperations, toward 0 on defections. At each turn, if w >= 0.5 *ATFT* cooperates; otherwise, it defects. The initial value are w = 0.5 and r=0.99 for the adaption rate.

**CCD** (*Periodically CCD*) -- Cooperates twice and then defects; then repeats the pattern.

**CD** (*Periodically CD*) -- Cooperates and then defects; then repeats the pattern.

**CS** (*Collective Strategy*) -- Starts cooperating and defects in the second turn. If the opponent played the same moves, *CS* starts behaving as *TFT*; otherwise, it plays AD.

**CTFT** (*Contrite Tit for Tat*) -- Same as *TFT* when there is no noise. In a noisy environment, once it receives T because of an error, it will choose cooperate twice in order to recover mutual cooperation.

**DDC** (*Periodically DDC*) -- Defects twice and then cooperates; then repeats the pattern.

**EXT2** (*Zero-Determinant Extort-2*) -- Let Szd be the total payoff earned by the ZD strategy, and So the one of its opponent, the strategy imposes the linear relationship Szd - **P** = 2 x (So - **P**). That is, guarantees to EXT2 twice the share of payoffs above the "punishment" threshold, compared to those received by the current opponent.

**FBF** (*Firm But Fair*) -- Cooperates on the first turn, and cooperates except after receiving a **S** payoff.

**FRT3** (*Fortress3*) -- Tries to recognize a kin member by playing the sequence *defect*, *defect*, *cooperate*. If the opponent plays the same sequence, it cooperates until the opponent defects. Otherwise, it defects until the opponent defects on continuous two moves, and then it cooperates on the following move.

**FRT4** (*Fortress4*) -- Tries to recognize the opponent by playing the sequence *defect*, *defect*, *defect*, *cooperate*. If the opponent plays the same sequence, it cooperates until the opponent defects. Otherwise, if it does not recognize the opponent as a friend, it defects until the opponent defects for three consecutive moves, then cooperates on the next.

**FS** (*Fair strategy*) -- Defects with a probability p equals to the frequency of defections played by the opponent; cooperates with probability 1-p. In the first turns the probability is smoothed towards 0.5.

**GRD** (*Gradual*) -- Cooperates on the first turn, and cooperates as long as the opponent cooperates; after the first defection of the other player, it defects one time and cooperates 2 times; ... after the *n*-th defection, it reacts with *n* consecutive defections and then plays two cooperate moves.

**GRM** (*Grim Trigger*) -- cooperates until the opponent defects, and subsequently always defects. The strategy is also known as *Grudger* or *Spiteful*.

**GTFT** (*Generous Tit for Tat*) -- Acts as a simple *TFT*, except that it cooperates with a probability q instead of just retaliating after the opponent defects. The parameter e is usually small, typically *e*=0.1. Some authors propose *e*=0.33 with the standard payoff matrix. The strategy is also known as *Naive Peace Maker* or *Soft Joss*.

**HM** (*Hard Majority*) -- Defects on the first turn, and defects if the number of defections of the opponent is greater than or equal to the number of times it has cooperated; otherwise it cooperates.

**HS** (*Handshake*) -- Defects on the first turn and cooperates on the second. If the opponent behaves the same, it always cooperates. Otherwise, it always defects.

**HTF2T** (*Hard Tit For Two Tats*) -- Cooperates unless the opponent plays two consecutive defections, then keeps defecting unless the adversary plays two consecutive cooperations. Then it starts cooperating again, and so on.

**HTFT2** (*Hard Tit for Tat (2-turn window)* -- Cooperates on the first turn, then defects only if the opponent has defected in any of the previous two turns. As some sources report the size of the window to be three turns instead of two, both versions have been included here.

**HTFT3** (*Hard Tit for Tat (3-turn window)* -- Cooperates on the first turn, then defects only if the opponent has defected in any of the previous three turns. See the discussion on the previous entry.

**NP** (*Naive Prober*) -- Acts as simple *TFT*, but defects unprovoked with a probability *e*. The parameter e is usually small, typically *e*=0.1. The strategy is also known as *Hard Joss*.

**OTFT** (*Omega Tit for Tat*) -- Normally behaves like *TFT*, but it is ready to play an extra cooperation for breaking deadlocks, i.e., two interlaced sequences of d alternating cooperation/defection. Moreover, *OTFT* measures the *randomness* of the opponent counting the number of times it changes behavior, and turns to an *Always Defect* if the value exceeds a given threshold t. Typically, *d*=3 and *t*=8.

**PAV** (*Pavlov*) -- Cooperates on the first turn; if a payoff of R or T is received in the last turn then repeats last choice, otherwise chooses the opposite one.

**PRO** (*Prober*) -- Starts with a sequence of one defection followed by two cooperations. If the opponent cooperated in the second and third turn, it keeps defecting for the rest of the game; otherwise, it plays as *TFT*.

**RND** (*Random Player*) -- Randomly chooses between cooperation and defection with equal probability, with no memory of previous exchanges.

**RP** (*Remorseful Prober*) -- Acts as simple *TFT*, but occasionally defects with a probability *e*. Unlike *Naive Prober*, however, it does not retort if the opponent defects after its unfair move. The parameter *e* is usually small, typically *e*=0.1.

**RTFT** (*Reverse Tit for Tat*) -- Starts defecting, and then chooses the opposite of the opponent's previous action. This apparently illogical variant of *Tit for Tat* is also known as *Psycho*.

**SGS** (*Southampton Group strategies*) -- A group of strategies are designed to recognize each other through a predetermined sequence of 5-10 moves at the start. Once two SGSs recognize each other, they will act as a *master* or *slave* - a master will always defect while a slave will always cooperate in order for the master to win the maximum points. If the opponent is not recognized as *SGS*, it will behave like an *AD* to minimize the score of the opponent.

**SG** (*Soft Grudger*) -- Cooperates until the opponent defects. In this case, it punishes the behavior with a sequence of four defections. Then, it offers a peace agreement with two consecutive moves of cooperation.

**SM** (*Soft Majority*) -- Cooperates on the first turn, and cooperates as long as the number of times the opponent has cooperated is greater than or equal to the number of times it has defected; otherwise, it defects.

**STF2T** (*Soft Tit For Two Tats*) -- Cooperates unless the opponent plays two consecutive defections.

**STFT** (*Suspicious Tit for Tat*) -- starts defecting, and then replicates the opponent's previous action. The strategy is also known as *Evil Tit For Tat*.

**TFT** (*Tit for Tat*) -- Cooperates on the first turn, copies the opponent's last move afterwards.

**ZDE** (*Zero-Determinant Extort*) -- let Szd be the total payoff earned by the ZD strategy, and So the one of its opponent, the strategy imposes the linear relationship Szd + **P** = 3 x (So - **P**).

**ZDF** (*Zero-Determinant Fixed Score*) -- let Szd be the total payoff earned by the ZD strategy, and So the one of its opponent, the strategy tries to fix the opponent's score to a pre-determined value g. Usually, with the standard payoff matrix, g=2.

**ZDG** (*Zero-Determinant Generous*) -- let Szd be the total payoff earned by the ZD strategy, and So the one of its opponent, the strategy enforces the relationship Szd = 2 x (So - **R**) between the two strategies' scores. Compared to EXT2, it offers the opponent a higher portion of the payoffs.


