# Protocols

Coming from a background in web development, I've always taken HTTP for granted. I could regurgitate facts about it: it's a protocol, a set of rules—easy peasy. And I could work with it just fine. That's the value of black boxing some complex topic—you can drive a car without knowing the specifics of internal combustion.

But now that I'm beginning to implement various protocols myself, I'm forced to confront the fact that I've never really given them much consideration beyond being aware of their existence.

I wanted to write down some bits of intuition I've hobbled together while journeying down this rabbit hole.

## Example: Two Prisoners

Let's say you have two prisoners who can neither see nor hear one another. The only thing connecting them is a bit of spare rope sprawled across the cement floor haphazardly.

If they wanted to communicate, they’d have to use this rope. Great. How are they going to do that?

### Initial Attempt

Suppose they’d agreed beforehand to swing the rope up and down to cause a wave. Then, using the alphabet to their advantage, they’d use a simple encoding where the letter 'a' was 1, and the letter 'z' was 26. Fine.

Well... not fine. They’re going to quickly run into a load of questions:

- How do you know when to stop counting the number of waves for each letter? For example, if you receive 3 waves, that could be: 'ab' or 'c'.
- How do you know when the message has stopped and there isn’t just a delay?
- What if the other person wants to reply? If they reply using the same rope, there will be collisions.
- What if the wind blows the rope a bit causing _random noise_, could you filter that out somehow?

### Need for Structure

It seems they need a more structured way to communicate here.

_Enter the protocol_

Defining the Basic Unit of Information
These problems generalize to any form of communication and any medium in which that communication occurs. One of the first things you may need to decide on is: what is the basic unit of information?

We said that these two prisoners used a single wave to represent a unit of communication. With this, you’d have to count each wave and assign it to some value.

It works, but that’s just SO many waves to count.

Wouldn't it be nicer if we didn’t have to send so many just to represent a single character?

### Introducing Multiple States

The problem is that our unit of communication can only represent one state. What if instead of a single wave of the rope, we just agreed that we could send one of two different types of waves: a _really big_ one and a _really small_ one? The emphasis on the word _really_ is important. If the difference between the two isn’t large enough, we may start guessing if we sent a small wave or a large one, thereby introducing _ambiguity_. We can't have that because now the prisoners would be stuck interpreting the amplitude of the waves.

Okay, so let's say that a wave with a height of a foot or more from bottom to top is big and one that is between 1 to 2 inches is small.

What does this afford us?

### Advantages Of Multiple States

With this new idea, we can tremendously cut down on the number of waves we'd need to send for our message. Let's still map the letters 'a' through 'z' in the same way we did before, with 'z' being 26. But let's say that we'll map a _big_ wave to the number 1 and the _small_ wave to the number 0. We... did we just reinvent binary? Yeah, I suppose we did.

One benefit of this is increased efficiency. Instead of sending 26 waves to represent the letter 'z', we can send this wave sequence:

```
big wave, big wave, small wave, big wave, small wave
```

Or in other terms:

```
11010
```

Well, in binary, this represents the number _26_ which is... the letter 'z'. So wait a minute. We only have to send 5 waves to represent something that used to take us 26 waves? Yep. That's the power of being able to represent multiple states.

That's awesome!

Now you may be tempted to say something like: well what about 3 states? Trinary? That would cut down the amount of waves we'd have to send for sure! And you'd be right.

But we'd ultimately introduce more ambiguity, because now, the prisoners may have to really be paying attention to the size of the waves coming across. This increased need for measurement can introduce bugs and complicate things. So we'll stick with binary for now.

### Data About Data

This is really great! The prisoners can now send less waves to communicate with one another. But there are still questions lingering.

- How does the other prisoner know when the message has stopped?
- How big is the message? Will the prisoner need to write it down because it's so large? This is useful information.
- Is the message still English or are there other things that can be communicated, like visuals? Pictures? Music?

I mean, look they probably just want to chit chat. They're not sending PDFs over a bit of spare rope. But let's stick with this example. Would it be possible to send information like that? Absolutely. Now whether or not the time would be reasonable for them is another story, but let's not worry about that for now.
